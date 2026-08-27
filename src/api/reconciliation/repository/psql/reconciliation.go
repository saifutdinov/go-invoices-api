package repository

import (
	"context"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (rr *ReconciliationRepository) CreatePaymentAllocation(ctx context.Context, pa domain.PaymenAllocation) error {
	_, err := rr.ExecContext(ctx, "INSERT INTO payment_allocation (payment_id, invoice_id, amount) VALUES ($1, $2, $3)", pa.PaymentId, pa.InvoiceId, pa.Amount)
	return err
}

func (rr *ReconciliationRepository) SumByInvoice(ctx context.Context, invoiceId uuid.UUID) (int64, error) {
	var sum int64
	if err := rr.QueryRowContext(ctx, "SELECT SUM(amount) FROM payment_allocation WHERE invoice_id = $1", invoiceId).Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}
