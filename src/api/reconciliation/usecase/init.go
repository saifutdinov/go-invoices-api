package usecase

import (
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/chronos"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
)

type ReconciliationUsecase struct {
	db.TXI
	LogRepository            domain.LogRepository
	ReconciliationRepository domain.ReconciliationReposotory
	InvoicesRepository       domain.InvoiceRepository
	PaymentsRepository       domain.PaymentRepository
	Scheduler                *chronos.Chronos
}

func NewReconciliationUsecase(
	db db.TXI,
	logRepo domain.LogRepository,
	recoRepo domain.ReconciliationReposotory,
	invoiceRepo domain.InvoiceRepository,
	paymentRepo domain.PaymentRepository,
	scheduler *chronos.Chronos,
) domain.ReconciliationUsecase {
	return &ReconciliationUsecase{
		TXI:                      db,
		LogRepository:            logRepo,
		ReconciliationRepository: recoRepo,
		InvoicesRepository:       invoiceRepo,
		PaymentsRepository:       paymentRepo,
		Scheduler:                scheduler,
	}
}
