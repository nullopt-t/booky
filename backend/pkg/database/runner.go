package database

import (
	"context"
)

type Runner interface {
	WithTx(ctx context.Context, fn func(tx QueryExecutor) error) error
	DB() QueryExecutor
}

type TxRunner struct {
	db *DB
}

func NewTxRunner(db *DB) *TxRunner {
	return &TxRunner{db: db}
}

func (t *TxRunner) WithTx(ctx context.Context, fn func(tx QueryExecutor) error) error {
	tx, err := t.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx) // safe rollback if commit never happens
	}()

	if err := fn(&Tx{tx}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (t *TxRunner) DB() QueryExecutor {
	return t.db
}
