package domain

import (
	"context"
	"time"
	"uuid"
)

type PaymenAllocation struct {
	PaymentId uuid.UUID
	InvoiceId uuid.UUID
	Amount    int64
}

type ReconciliationReport struct {
	Invoices []InvoiceReport     `json:"invoices"`
	Payments []PaymentReportData `json:"payments"`
}

type InvoiceReport struct {
	ID              uuid.UUID     `json:"id"`
	ExternalID      string        `json:"reference"`
	Amount          int64         `json:"amount"`
	PaidAmount      int64         `json:"paid_amount"`
	RemainingAmount int64         `json:"remaining_amount"`
	Status          InvoiceStatus `json:"status"`
	DueDate         time.Time     `json:"due_date"`
	PaymentsCount   int           `json:"payments_count"`
}

type PaymentReportData struct {
	ID          uuid.UUID     `json:"id"`
	Amount      int64         `json:"amount"`
	Reference   string        `json:"reference"`
	PaymentDate time.Time     `json:"payment_date"`
	Status      PaymentStatus `json:"status"`

	InvoiceID       *uuid.UUID `json:"invoice_id"`
	InvoiceAmount   *int64     `json:"invoice_amount"`
	AllocatedAmount int64      `json:"allocated_amount"`

	Discrepancy *PaymentDiscrepancy `json:"discrepancy,omitempty"`
}

type PaymentDiscrepancy struct {
	Type           string `json:"type"`
	Message        string `json:"message"`
	ExpectedAmount int64  `json:"expected_amount,omitempty"`
	ActualAmount   int64  `json:"actual_amount,omitempty"`
	Difference     int64  `json:"difference,omitempty"`
}

type ReconciliationReposotory interface {
	CreatePaymentAllocation(ctx context.Context, reconciliation PaymenAllocation) error
	SumByInvoice(ctx context.Context, invoiceId uuid.UUID) (int64, error)
}

type ReconciliationUsecase interface {
	StartBackgroundProcessing()

	ProcessPayments(ctx context.Context)

	Report(ctx context.Context) (*ReconciliationReport, error)
}
