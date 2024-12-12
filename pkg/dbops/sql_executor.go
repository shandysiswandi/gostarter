// Package dbops provides abstractions and utility functions for database operations.
// It defines interfaces for querying and executing SQL commands, and includes helper functions
// for performing common database operations, such as executing queries and scanning rows.
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

func prepare(qp QueryProvider) (string, []any, error) {
	q, args, err := qp()
	if err != nil {
		if verbose {
			log.Printf("Exec:queryProvider query: %s, args: %v, err: %v \n", q, args, err)
		}

		return "", nil, err
	}

	return q, args, nil
}

// Exec executes a query and handles the result.
func Exec(ctx context.Context, exec Execer, qp QueryProvider, feedback ...bool) error {
	query, args, err := prepare(qp)
	if err != nil {
		return err
	}

	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		exec = tx
	}

	res, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		if verbose {
			log.Printf("Exec:ExecContext result: %v, query: %s, err: %v \n", res, query, err)
		}

		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		if verbose {
			log.Printf("Exec:RowsAffected rows: %d, query: %s, err: %v \n", aff, query, err)
		}

		return err
	}

	if aff == 0 && len(feedback) > 0 && feedback[0] {
		if verbose {
			log.Printf("Exec: query: %s, data not affected or not found \n", query)
		}

		return ErrZeroRowsAffected
	}

	return nil
}

// SQLGet executes a query that is expected to return a single row and scans the result into a type T.
// It uses the provided Row[T] implementation to handle scanning of columns.
func SQLGet[T any, PT Row[T]](ctx context.Context, q Queryer, qp QueryProvider) (*T, error) {
	query, args, err := prepare(qp)
	if err != nil {
		return nil, err
	}

	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		q = tx
	}

	var t T
	ptr := PT(&t)

	err = q.QueryRowContext(ctx, query, args...).Scan(ptr.ScanColumn()...)
	if errors.Is(err, sql.ErrNoRows) {
		if verbose {
			log.Printf("Exec: query: %s, data not found \n", query)
		}

		return nil, nil //nolint:nilnil // no rows is not an error, just a nil result
	}

	if err != nil {
		if verbose {
			log.Printf("Exec:QueryRowContext query: %s, err: %v \n", query, err)
		}

		return nil, err
	}

	return &t, nil
}

// SQLGets executes a query that may return multiple rows and scans the results into a slice of type T.
// It uses the provided Row[T] implementation to handle scanning of columns.
func SQLGets[T any, PT Row[T]](ctx context.Context, q Queryer, qp QueryProvider) ([]T, error) {
	query, args, err := prepare(qp)
	if err != nil {
		return nil, err
	}

	if verbose {
		log.Printf("SQLGets: query: %s\n", query)
	}

	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if ok {
		q = tx
	}

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		if verbose {
			log.Printf("Exec:QueryRowContext query: %s, err: %v \n", query, err)
		}

		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v\n", err)
		}
	}()

	var entities []T
	for rows.Next() {
		var t T
		ptr := PT(&t)

		if err := rows.Scan(ptr.ScanColumn()...); err != nil {
			if verbose {
				log.Printf("Exec:ScanColumn query: %s, err: %v \n", query, err)
			}

			return nil, ErrScanRow
		}

		entities = append(entities, t)
	}

	if err := rows.Err(); err != nil {
		if verbose {
			log.Printf("Exec:rows.Err query: %s, err: %v \n", query, err)
		}

		return nil, err
	}

	return entities, nil
}
