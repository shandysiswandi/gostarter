package sqlkit

import "context"

type NoopDB struct{}

func NewNoopDB() *NoopDB {
	return &NoopDB{}
}

func (*NoopDB) Transaction(ctx context.Context, fn func(ctx context.Context) error, _ ...TxOption) error {
	return fn(ctx)
}
