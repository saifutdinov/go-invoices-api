package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/chronos"
)

const (
	ReconciliationInterval = time.Second * 20
)

func (ru *ReconciliationUsecase) StartBackgroundProcessing() {
	log.Println("Chronos: background tasks initiated")
	ru.scheduleNextReconciliation(time.Now().Add(ReconciliationInterval))
}

func (ru *ReconciliationUsecase) scheduleNextReconciliation(startAt time.Time) {
	task := chronos.NewTask(
		startAt,
		func(ctx context.Context) {
			ru.ProcessPayments(ctx)
		},
	)

	ru.Scheduler.Schedule(task)
}

func (ru *ReconciliationUsecase) ProcessPayments(ctx context.Context) {
	defer ru.scheduleNextReconciliation(time.Now().Add(ReconciliationInterval))

	pendingPayments, err := ru.PaymentsRepository.ReadPending(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, paymentID := range pendingPayments {
		if err := ru.reconcile(ctx, paymentID); err != nil {
			log.Println(err)

			ru.LogRepository.Create(
				ctx,
				domain.Log{
					EventType:  domain.LogEventReconciliationFailed,
					EntityType: "payment",
					EntityID:   &paymentID,
					Message:    "Payment reconciliation failed",
					Metadata: map[string]any{
						"error": err.Error(),
					},
				},
			)

			continue
		}
	}
}

func (ru *ReconciliationUsecase) reconcile(ctx context.Context, paymentId uuid.UUID) error {
	txCtx, err := ru.BeginTx(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	defer ru.RollbackTx(txCtx)

	payment, err := ru.PaymentsRepository.ReadForUpdate(txCtx, paymentId)
	if err != nil {
		log.Println(err)
		return err
	}

	if payment.Status != domain.PaymentStatusPending {

		ru.LogRepository.Create(
			txCtx,
			domain.Log{
				EventType:  domain.LogEventPaymentAlreadyHandled,
				EntityType: "payment",
				EntityID:   &payment.ID,
				Message:    "Payment was already processed",
			},
		)

		return nil
	}

	invoice, err := ru.InvoicesRepository.ReadForUpdate(txCtx, payment.Reference)
	if err != nil {
		log.Println(err)

		if errors.Is(err, sql.ErrNoRows) {
			ru.LogRepository.Create(
				txCtx,
				domain.Log{
					EventType:  domain.LogEventPaymentUnmatched,
					EntityType: "payment",
					EntityID:   &payment.ID,
					Message:    "No matching invoice found",
					Metadata: map[string]any{
						"reference":    payment.Reference,
						"amount":       payment.Amount,
						"payment_date": payment.PaymentDate,
					},
				},
			)
		}

		return err
	}

	ru.LogRepository.Create(
		txCtx,
		domain.Log{
			EventType:  domain.LogEventInvoiceMatched,
			EntityType: "invoice",
			EntityID:   &invoice.ID,
			Message:    "Payment matched with invoice",
			Metadata: map[string]any{
				"payment_id": payment.ID,
				"reference":  payment.Reference,
			},
		},
	)

	if err := ru.ReconciliationRepository.CreatePaymentAllocation(txCtx, domain.PaymenAllocation{
		PaymentId: payment.ID,
		InvoiceId: invoice.ID,
		Amount:    payment.Amount,
	}); err != nil {
		log.Println(err)
		return err
	}

	ru.LogRepository.Create(
		txCtx,
		domain.Log{
			EventType:  domain.LogEventPaymentAllocated,
			EntityType: "payment",
			EntityID:   &payment.ID,
			Message:    "Payment allocated to invoice",
			Metadata: map[string]any{
				"invoice_id": invoice.ID,
				"amount":     payment.Amount,
			},
		},
	)

	// Calculate total amount paid for invoice.
	paidAmount, err := ru.ReconciliationRepository.SumByInvoice(
		txCtx,
		invoice.ID,
	)
	if err != nil {
		return err
	}

	invoiceStatus := calculateInvoiceStatus(
		invoice.Amount,
		paidAmount,
	)

	if invoice.Status != invoiceStatus {

		if err := ru.InvoicesRepository.UpdateStatus(txCtx, invoice.ID, invoiceStatus); err != nil {
			log.Println(err)
			return err
		}

		ru.LogRepository.Create(
			txCtx,
			domain.Log{
				EventType:  domain.LogEventInvoiceStatusChanged,
				EntityType: "invoice",
				EntityID:   &invoice.ID,
				Message:    "Invoice status changed",
				Metadata: map[string]any{
					"old_status": invoice.Status,
					"new_status": invoiceStatus,
				},
			},
		)

	}

	if err := ru.PaymentsRepository.UpdateStatus(txCtx, payment.ID, domain.PaymentStatusCompleted); err != nil {
		log.Println(err)
		return err
	}

	if err := ru.CommitTx(txCtx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func calculateInvoiceStatus(
	invoiceAmount int64,
	paidAmount int64,
) domain.InvoiceStatus {
	switch {
	case paidAmount == 0:
		return domain.InvoiceStatusUnmatched
	case paidAmount < invoiceAmount:
		return domain.InvoiceStatusPartiallyPaid
	case paidAmount == invoiceAmount:
		return domain.InvoiceStatusPaid
	default:
		return domain.InvoiceStatusOverpaid
	}
}

func (ru *ReconciliationUsecase) Report(ctx context.Context) (*domain.ReconciliationReport, error) {

	invoices, err := ru.InvoicesRepository.ReadReport(ctx)
	if err != nil {
		return nil, err
	}

	payments, err := ru.PaymentsRepository.ReadReport(ctx)
	if err != nil {
		return nil, err
	}

	for idx, payment := range payments {
		payments[idx].Discrepancy = buildPaymentDiscrepancy(payment)
	}

	return &domain.ReconciliationReport{
		Invoices: invoices,
		Payments: payments,
	}, nil
}

func buildPaymentDiscrepancy(
	payment domain.PaymentReportData,
) *domain.PaymentDiscrepancy {

	if payment.InvoiceID == nil {
		return &domain.PaymentDiscrepancy{
			Type:    "invoice_not_found",
			Message: "No invoice found for payment reference",
		}
	}

	if payment.AllocatedAmount < *payment.InvoiceAmount {
		return &domain.PaymentDiscrepancy{
			Type:           "underpaid",
			Message:        "Invoice is partially paid",
			ExpectedAmount: *payment.InvoiceAmount,
			ActualAmount:   payment.AllocatedAmount,
			Difference:     payment.AllocatedAmount - *payment.InvoiceAmount,
		}
	}

	if payment.AllocatedAmount > *payment.InvoiceAmount {
		return &domain.PaymentDiscrepancy{
			Type:           "overpaid",
			Message:        "Invoice is overpaid",
			ExpectedAmount: *payment.InvoiceAmount,
			ActualAmount:   payment.AllocatedAmount,
			Difference:     payment.AllocatedAmount - *payment.InvoiceAmount,
		}
	}

	return nil
}
