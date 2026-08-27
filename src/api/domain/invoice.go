package domain

import (
	"context"
	"strings"
	"uuid"
)

type InvoiceStatus string

// Paid, Partially Paid, Overpaid, Unmatched
const (
	InvoiceStatusUnmatched     InvoiceStatus = "unmatched"
	InvoiceStatusPartiallyPaid InvoiceStatus = "partially paid"
	InvoiceStatusOverpaid      InvoiceStatus = "overpaid"
	InvoiceStatusPaid          InvoiceStatus = "paid"
)

func (is *InvoiceStatus) UnmarshalJSON(b []byte) error {
	v := strings.ToLower(string(b))
	*is = InvoiceStatus(strings.Trim(v, `"`))
	return nil
}

type Invoice struct {
	ID         uuid.UUID
	ExternalID string
	Amount     int64
	DueDate    string
	CreatedAt  int64
	UpdatedAt  int64
	Status     InvoiceStatus
}

type InvoiceUsecase interface {
	Save(ctx context.Context, invoice Invoice) error
}

type InvoiceRepository interface {
	Save(ctx context.Context, invoice Invoice) error
	ReadForUpdate(ctx context.Context, invoiceExternalId string) (Invoice, error)
	UpdateStatus(ctx context.Context, invoiceId uuid.UUID, status InvoiceStatus) error

	ReadReport(ctx context.Context) ([]InvoiceReport, error)
}
