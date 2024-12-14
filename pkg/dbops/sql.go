// Package dbops provides abstractions and utility functions for database operations.
// It defines interfaces for querying and executing SQL commands, and includes helper functions
// for performing common database operations, such as executing queries and scanning rows.
package dbops

import (
	"context"
	"database/sql"
)

// Queryer defines methods for querying the database.
type Queryer interface {
	// QueryRowContext executes a query that is expected to return a single row.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	// QueryContext executes a query that may return multiple rows.
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
}

// Execer defines methods for executing SQL commands.
type Execer interface {
	// ExecContext executes a query without returning any rows.
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

// Row represents a single row from a database query with scanning capability.
type Row[T any] interface {
	// ScanColumn returns a slice of values scanned from the row.
	ScanColumn() []any

	// Embed the value of type T.
	*T
}

// Tx defines the interface for managing SQL transactions.
type Tx interface {
	// Transaction executes the provided function within a transaction context.
	// It starts a transaction, passes a context containing the transaction,
	// and ensures that the transaction is committed or rolled back based on success or failure.
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type SQLTx interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

type QueryProvider func() (string, []any, error)
