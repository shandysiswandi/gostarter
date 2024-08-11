package persistence

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

// investigate: https://github.com/wroge/scan

var ErrScanRow = errors.New("failed to scan column into field type")

func SQLGet[T any, PT Row[T]](
	ctx context.Context,
	querier Queryer,
	queryProvider func() (string, []any, error),
) (*T, error) {
	query, args, err := queryProvider()
	if err != nil {
		return nil, err
	}

	var t T
	ptr := PT(&t)
	err = querier.QueryRowContext(ctx, query, args...).Scan(ptr.Columns()...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil //nolint:nilnil // no rows is not an error, just a nil result
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func SQLGets[T any, PT Row[T]](
	ctx context.Context,
	querier Queryer,
	queryProvider func() (string, []any, error),
) ([]T, error) {
	query, args, err := queryProvider()
	if err != nil {
		return nil, err
	}

	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
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

		if err := rows.Scan(ptr.Columns()...); err != nil {
			return nil, ErrScanRow
		}

		entities = append(entities, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return entities, nil
}
