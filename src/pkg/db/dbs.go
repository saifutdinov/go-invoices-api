package db

import (
	"context"
	"database/sql"
)

type DBS struct {
	dbptr *sql.DB
}

func NewDBS(dbp *sql.DB) DBI {
	return &DBS{dbptr: dbp}
}

func (dbs *DBS) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return switchDB(ctx, dbs.dbptr).ExecContext(ctx, query, args...)
}

func (dbs *DBS) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return switchDB(ctx, dbs.dbptr).PrepareContext(ctx, query)
}

func (dbs *DBS) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return switchDB(ctx, dbs.dbptr).QueryContext(ctx, query, args...)
}

func (dbs *DBS) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return switchDB(ctx, dbs.dbptr).QueryRowContext(ctx, query, args...)
}
