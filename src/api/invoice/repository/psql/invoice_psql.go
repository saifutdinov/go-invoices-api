package psql

import (
	"context"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (ir *InvoiceRepository) Save(ctx context.Context, invoice domain.Invoice) error {
	invoice.ID = uuid.New()
	_, err := ir.ExecContext(ctx, "INSERT INTO invoice (id, external_id, amount, due_date, status) VALUES ($1, $2, $3, $4, $5)", invoice.ID, invoice.ExternalID, invoice.Amount, invoice.DueDate, invoice.Status)
	return err
}

func (ir *InvoiceRepository) ReadForUpdate(ctx context.Context, invoiceExternalId string) (domain.Invoice, error) {
	var invoice domain.Invoice
	if err := ir.QueryRowContext(ctx, "SELECT id, external_id, amount, due_date, status FROM invoice WHERE external_id = $1 FOR UPDATE", invoiceExternalId).Scan(
		&invoice.ID,
		&invoice.ExternalID,
		&invoice.Amount,
		&invoice.DueDate,
		&invoice.Status,
	); err != nil {
		return domain.Invoice{}, err
	}
	return invoice, nil
}

func (ir *InvoiceRepository) UpdateStatus(ctx context.Context, invoiceId uuid.UUID, status domain.InvoiceStatus) error {
	_, err := ir.ExecContext(ctx, "UPDATE invoice SET status = $1, updatedat = now() WHERE id = $2", status, invoiceId)
	return err
}

func (ir *InvoiceRepository) ReadReport(ctx context.Context) ([]domain.InvoiceReport, error) {

	rows, err := ir.QueryContext(ctx, `
	SELECT
		i.id,
		i.external_id,
		i.amount,
		i.status,
		i.due_date,

		COALESCE(SUM(pa.amount), 0) AS paid_amount,
		COUNT(pa.payment_id) AS payments_count,

		i.amount - COALESCE(SUM(pa.amount), 0) AS remaining_amount

	FROM invoice i

	LEFT JOIN payment_allocation pa
		ON pa.invoice_id = i.id::varchar
	GROUP BY
		i.id,
		i.external_id,
		i.amount,
		i.status,
		i.due_date

	ORDER BY i.due_date ASC;
	`)
	if err != nil {
		return nil, err
	}

	var reports []domain.InvoiceReport
	for rows.Next() {
		var report domain.InvoiceReport
		if err := rows.Scan(
			&report.ID,
			&report.ExternalID,
			&report.Amount,
			&report.Status,
			&report.DueDate,

			&report.PaidAmount,
			&report.PaymentsCount,

			&report.RemainingAmount,
		); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil

}
