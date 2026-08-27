package db

import (
	"context"
	"database/sql"
)

// sqlTxKeyType type for sql transaction context key
type sqlTxKeyType string

// sqlTxValueKey key for sql transaction stored in child context
const sqlTxValueKey = sqlTxKeyType("sql_tx_value_key")

type TXI interface {
	// BeginTx creates child context that stores started transaction with sqlTxValueKey as key
	BeginTx(ctx context.Context) (context.Context, error)
	// CommitTx commits transaction in created child context
	CommitTx(ctx context.Context) error
	// RollbackTx rollbacks transaction in created child context
	RollbackTx(ctx context.Context) error
}

type DBI interface {
	// ExecContext method from sql.DB and sql.Tx both
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// PrepareContext method from sql.DB and sql.Tx both
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// QueryContext method from sql.DB and sql.Tx both
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// QueryRowContext method from sql.DB and sql.Tx both
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
