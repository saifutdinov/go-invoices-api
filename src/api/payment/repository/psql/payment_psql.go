package psql

import (
	"context"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (pr *PaymentRepository) Save(ctx context.Context, payment domain.Payment) error {
	payment.ID = uuid.New()
	_, err := pr.ExecContext(ctx, "INSERT INTO payment (id, external_id, amount, payment_date, reference, source, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		payment.ID,
		payment.ExternalID,
		payment.Amount,
		payment.PaymentDate,
		payment.Reference,
		payment.Source,
		payment.Status,
	)
	return err
}

func (pr *PaymentRepository) ReadPending(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := pr.QueryContext(ctx, "SELECT id FROM payment WHERE status = 'pending' LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []uuid.UUID
	for rows.Next() {
		var payment uuid.UUID
		if err := rows.Scan(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (pr *PaymentRepository) ReadForUpdate(ctx context.Context, paymentId uuid.UUID) (domain.Payment, error) {
	var payment domain.Payment
	if err := pr.QueryRowContext(ctx, "SELECT id, external_id, amount, payment_date, reference, source, status FROM payment WHERE id = $1 FOR UPDATE", paymentId).Scan(
		&payment.ID,
		&payment.ExternalID,
		&payment.Amount,
		&payment.PaymentDate,
		&payment.Reference,
		&payment.Source,
		&payment.Status,
	); err != nil {
		return domain.Payment{}, err
	}
	return payment, nil
}

func (pr *PaymentRepository) UpdateStatus(ctx context.Context, paymentId uuid.UUID, status domain.PaymentStatus) error {
	_, err := pr.ExecContext(ctx, "UPDATE payment SET status = $2, updatedat = now() WHERE id = $1", paymentId, status)
	return err
}

func (pr *PaymentRepository) ReadReport(ctx context.Context) ([]domain.PaymentReportData, error) {

	rows, err := pr.QueryContext(ctx, `
	SELECT
		p.id,
		p.amount,
		p.reference,
		p.payment_date,
		p.status,

		i.id AS invoice_id,
		i.amount AS invoice_amount,

		COALESCE(SUM(pa.amount), 0) AS allocated_amount

	FROM payment p

	LEFT JOIN invoice i
		ON i.external_id = p.reference

	LEFT JOIN payment_allocation pa
		ON pa.payment_id = p.id::varchar

	GROUP BY
		p.id,
		p.amount,
		p.reference,
		p.payment_date,
		p.status,
		i.id,
		i.amount

	ORDER BY p.payment_date DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []domain.PaymentReportData
	for rows.Next() {
		var payment domain.PaymentReportData
		if err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Reference,
			&payment.PaymentDate,
			&payment.Status,
			&payment.InvoiceID,
			&payment.InvoiceAmount,
			&payment.AllocatedAmount,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}
