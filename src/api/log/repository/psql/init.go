package psql

import (
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
)

type LogRepository struct {
	db.DBI
}

func NewLogRepository(dbo db.DBI) domain.LogRepository {
	return &LogRepository{dbo}
}
