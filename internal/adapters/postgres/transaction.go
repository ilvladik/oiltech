package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (m *TransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return fn(ctx)
	}
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, txKey{}, tx)
	defer tx.Rollback()
	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

func TxOrDb(ctx context.Context, db *sqlx.DB) DBTX {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return db
}
