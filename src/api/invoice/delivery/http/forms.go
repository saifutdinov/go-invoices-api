package http

import (
	"encoding/json"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

type CreateInvoiceForm struct {
	ID      string               `json:"id"`
	Amount  json.Number          `json:"amount"`
	DueDate string               `json:"due_date"`
	Status  domain.InvoiceStatus `json:"status"`
}
