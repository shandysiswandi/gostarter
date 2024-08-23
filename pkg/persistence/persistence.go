// Package persistence provides abstractions and utility functions for database operations.
//
// This package defines interfaces for querying and executing SQL commands,
// and includes helper functions to simplify common database operations.
package persistence

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
