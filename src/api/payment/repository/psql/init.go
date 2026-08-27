package psql

import (
	"database/sql"

	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
)

type PaymentRepository struct {
	db.DBI
}

func NewPaymentRepository(dbo *sql.DB) domain.PaymentRepository {
	return &PaymentRepository{
		DBI: db.NewDBS(dbo),
	}
}
