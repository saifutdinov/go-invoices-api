package domain

import (
	"context"
	"uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
)

type Payment struct {
	ID          uuid.UUID
	ExternalID  string
	Amount      int64
	PaymentDate string
	Reference   string
	Source      string
	Status      PaymentStatus
}

type PaymentUsecase interface {
	Save(ctx context.Context, payment Payment) error
	UpdateStatus(ctx context.Context, paymentId uuid.UUID, status PaymentStatus) error
}

type PaymentRepository interface {
	Save(ctx context.Context, payment Payment) error
	ReadPending(ctx context.Context) ([]uuid.UUID, error)
	ReadForUpdate(ctx context.Context, paymentId uuid.UUID) (Payment, error)
	UpdateStatus(ctx context.Context, paymentId uuid.UUID, status PaymentStatus) error

	ReadReport(ctx context.Context) ([]PaymentReportData, error)
}
