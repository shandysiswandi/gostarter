package sqlkit

import (
	"context"
	"database/sql"
)

type Tx interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error
}

type contextKeySQLTx struct{}

type TxOption func(*sql.TxOptions)

func WithIsolationLevel(l sql.IsolationLevel) TxOption {
	return func(to *sql.TxOptions) {
		to.Isolation = l
	}
}

func WithReadOnly(enable bool) TxOption {
	return func(to *sql.TxOptions) {
		to.ReadOnly = enable
	}
}

func (d *DB) Transaction(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error {
	txOpts := &sql.TxOptions{}

	for _, opt := range opts {
		opt(txOpts)
	}

	tx, err := d.db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, contextKeySQLTx{}, tx)

	defer func() {
		if tx == nil {
			return
		}

		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				d.log.Error(ctx, "panic when execute function and error rollback", err)
			}

			d.log.Warn(ctx, "recover from panic when execute function")

			panic(r) // is re-panic
		}

		if err != nil {
			if err := tx.Rollback(); err != nil {
				d.log.Error(ctx, "error when execute function", err)
			}

			return
		}

		if err := tx.Commit(); err != nil {
			d.log.Error(ctx, "error when commit transaction", err)
		}
	}()

	err = fn(ctx)

	return err
}
