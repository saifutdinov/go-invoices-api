package http

import (
	"encoding/json"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

type CreatePaymentForm struct {
	ID          string      `json:"id" validate:"required"`
	Amount      json.Number `json:"amount" validate:"required"`
	PaymentDate string      `json:"payment_date" validate:"required"`
	Reference   string      `json:"reference" validate:"required"`
	Source      string      `json:"source" validate:"required"`
}

type UpdatePaymentStatusForm struct {
	Status domain.PaymentStatus `json:"status" validate:"required"`
}
