package dbops

import (
	"context"
	"database/sql"
)

type NoopDB struct{}

func SetContextNoopTx(ctx context.Context) context.Context {
	var tx *sql.Tx
	return context.WithValue(ctx, contextKeySQLTx{}, tx)
}

func NewNoopDB() *NoopDB {
	return &NoopDB{}
}

func (*NoopDB) BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error) {
	return nil, nil
}
