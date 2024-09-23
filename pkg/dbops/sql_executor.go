// Package dbops provides abstractions and utility functions for database operations.
// It defines interfaces for querying and executing SQL commands,
// and includes helper functions for performing common database operations, such as executing queries and scanning rows.
package dbops

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

var (
	// ErrZeroRowsAffected is returned when an update, insert, or delete operation affects zero rows.
	ErrZeroRowsAffected = errors.New("no rows affected by an update, insert, or delete")

	// ErrScanRow is returned when scanning a row into the field type fails.
	ErrScanRow = errors.New("failed to scan column into field type")
)

// Exec executes a query and handles the result.
func Exec(ctx context.Context, execer Execer, queryProvider func() (string, []any, error), feedback ...bool) error {
	// Generate the query and arguments from the queryProvider function.
	query, args, err := queryProvider()
	if err != nil {
		return err
	}

	// Check if the context contains an active transaction.
	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		execer = tx
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

// SQLGet executes a query that is expected to return a single row and scans the result into a type T.
// It uses the provided Row[T] implementation to handle scanning of columns.
func SQLGet[T any, PT Row[T]](
	ctx context.Context,
	querier Queryer,
	queryProvider func() (string, []any, error),
) (*T, error) {
	// Generate the query and arguments from the queryProvider function.
	query, args, err := queryProvider()
	if err != nil {
		return nil, err
	}

	// Check if the context contains an active transaction.
	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		querier = tx
	}

	// Initialize a variable of type T.
	var t T
	ptr := PT(&t)

	// Execute the query and scan the result into the variable.
	err = querier.QueryRowContext(ctx, query, args...).Scan(ptr.ScanColumn()...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil //nolint:nilnil // no rows is not an error, just a nil result
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// SQLGets executes a query that may return multiple rows and scans the results into a slice of type T.
// It uses the provided Row[T] implementation to handle scanning of columns.
func SQLGets[T any, PT Row[T]](
	ctx context.Context,
	querier Queryer,
	queryProvider func() (string, []any, error),
) ([]T, error) {
	// Generate the query and arguments from the queryProvider function.
	query, args, err := queryProvider()
	if err != nil {
		return nil, err
	}

	// Check if the context contains an active transaction.
	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		querier = tx
	}

	// Execute the query and obtain a result set.
	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}()

	// Iterate through the rows and scan each row into a slice of type T.
	var entities []T
	for rows.Next() {
		var t T
		ptr := PT(&t)

		if err := rows.Scan(ptr.ScanColumn()...); err != nil {
			return nil, ErrScanRow
		}

		entities = append(entities, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
