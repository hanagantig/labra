package mysqlxrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const txKey contextKey = "sql_tx"

type Transactor struct {
	conn *sqlx.DB
}

type conn interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func NewTransactor(c *sqlx.DB) Transactor {
	return Transactor{
		conn: c,
	}
}

func (t *Transactor) InTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := t.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	ctxCopy := NewCtxWithTrx(ctx, tx)

	err = txFunc(ctxCopy)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("%v:%w", rollbackErr, err)
		}

		return err
	}

	return tx.Commit()
}

func NewCtxWithTrx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func (t *Transactor) GetConn(ctx context.Context) conn {
	db, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		return t.conn
	}

	return db
}
