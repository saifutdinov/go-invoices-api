package usecase

import (
	"context"
	"uuid"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (pu *PaymentUsecase) Save(ctx context.Context, payment domain.Payment) error {
	if err := pu.PaymentRepository.Save(ctx, payment); err != nil {
		return err
	}

	pu.LogRepository.Create(
		ctx,
		domain.Log{
			EventType:  domain.LogEventPaymentReceived,
			EntityType: "payment",
			EntityID:   &payment.ID,
			Message:    "Payment received",
			Metadata: map[string]any{
				"amount":       payment.Amount,
				"reference":    payment.Reference,
				"payment_date": payment.PaymentDate,
				"source":       payment.Source,
			},
		},
	)

	return nil
}

func (pu *PaymentUsecase) UpdateStatus(ctx context.Context, paymentId uuid.UUID, status domain.PaymentStatus) error {
	return pu.PaymentRepository.UpdateStatus(ctx, paymentId, status)
}
