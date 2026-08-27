package repository

import (
	"database/sql"

	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
)

type ReconciliationRepository struct {
	db.DBI
}

func NewReconciliationRepository(dbo *sql.DB) domain.ReconciliationReposotory {
	return &ReconciliationRepository{
		DBI: db.NewDBS(dbo),
	}
}
