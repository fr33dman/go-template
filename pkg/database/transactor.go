package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlcrepo "github.com/fr33dman/go-template/internal/repo/entity/sqlc"
)

type txContextKey struct{}

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{pool: pool}
}

func (t *Transactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = fn(context.WithValue(ctx, txContextKey{}, tx))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func DBTX(ctx context.Context, fallback sqlcrepo.DBTX) sqlcrepo.DBTX {
	if tx, ok := Tx(ctx); ok {
		return tx
	}
	return fallback
}

func Tx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx, ok
}
