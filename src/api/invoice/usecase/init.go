package usecase

import (
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

type InvoiceUsecase struct {
	db.TXI
	Config            *toml.Config
	LogRepository     domain.LogRepository
	InvoiceRepository domain.InvoiceRepository
}

func NewInvoiceUsecase(
	tx db.TXI,
	logRepo domain.LogRepository,
	invoiceRepo domain.InvoiceRepository,
	config *toml.Config,
) domain.InvoiceUsecase {
	return &InvoiceUsecase{
		TXI:               tx,
		LogRepository:     logRepo,
		InvoiceRepository: invoiceRepo,
		Config:            config,
	}
}
