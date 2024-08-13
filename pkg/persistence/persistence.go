package persistence

import (
	"context"
	"database/sql"
)

type Queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
}

type Execer interface {
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

type Row[T any] interface {
	ScanColumn() []any
	*T
}
