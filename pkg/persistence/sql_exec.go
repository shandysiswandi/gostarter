// Package persistence provides abstractions and utility functions for database operations.
//
// This package defines interfaces for querying and executing SQL commands,
// and includes helper functions to simplify common database operations.
package persistence

import (
	"context"
	"errors"
)

// ErrZeroRowsAffected is returned when an update, insert, or delete operation affects zero rows.
var ErrZeroRowsAffected = errors.New("no rows affected by an update, insert, or delete")

// Exec executes a query and handles the result.
func Exec(ctx context.Context, execer Execer, queryProvider func() (string, []any, error), feedback ...bool) error {
	// Generate the query and arguments from the queryProvider function.
	query, args, err := queryProvider()
	if err != nil {
		return err
	}

	// Execute the query using the provided Execer.
	res, err := execer.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// Check the number of rows affected if feedback is requested.
	aff, err := res.RowsAffected()
	if err != nil && len(feedback) > 0 && feedback[0] {
		return err
	}

	if aff == 0 && len(feedback) > 0 && feedback[0] {
		return ErrZeroRowsAffected
	}

	return nil
}
