package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TXS struct {
	dbptr *sql.DB
}

func NewTransaction(dbp *sql.DB) TXI {
	return &TXS{dbptr: dbp}
}

func (txs *TXS) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := txs.dbptr.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, sqlTxValueKey, tx), nil
}

func (txs *TXS) CommitTx(ctx context.Context) error {
	if tx := ctx.Value(sqlTxValueKey); tx != nil {
		if err := tx.(*sql.Tx).Commit(); err != nil && err != sql.ErrTxDone {
			return err
		}
	}
	return nil
}

func (txs *TXS) RollbackTx(ctx context.Context) error {
	if tx := ctx.Value(sqlTxValueKey); tx != nil {
		if err := tx.(*sql.Tx).Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Println("Transaction, RollbackTx: ", err)
			return err
		}
	}
	return nil
}

// switchDB return created sql.Tx in child context if exists and sql.DB otherwise
func switchDB(ctx context.Context, dbp *sql.DB) DBI {
	if tx := ctx.Value(sqlTxValueKey); tx != nil {
		return tx.(*sql.Tx)
	} else {
		return dbp
	}
}
