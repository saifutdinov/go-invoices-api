package psql

import (
	"database/sql"

	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
)

type InvoiceRepository struct {
	db.DBI
}

func NewInvoiceRepository(dbo *sql.DB) domain.InvoiceRepository {
	return &InvoiceRepository{
		DBI: db.NewDBS(dbo),
	}
}
