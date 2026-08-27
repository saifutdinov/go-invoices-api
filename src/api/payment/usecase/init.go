package usecase

import (
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

type PaymentUsecase struct {
	db.TXI
	Config            *toml.Config
	LogRepository     domain.LogRepository
	PaymentRepository domain.PaymentRepository
}

func NewPaymentUsecase(
	tx db.TXI,
	logRepo domain.LogRepository,
	paymentRepo domain.PaymentRepository,
	config *toml.Config,
) domain.PaymentUsecase {
	return &PaymentUsecase{
		TXI:               tx,
		LogRepository:     logRepo,
		PaymentRepository: paymentRepo,
		Config:            config,
	}
}
