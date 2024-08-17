package persistence

import (
	"context"
	"errors"
)

var ErrZeroRowsAffected = errors.New("no rows affected by an update, insert, or delete")

func Exec(ctx context.Context, execer Execer, queryProvider func() (string, []any, error), feedback ...bool) error {
	query, args, err := queryProvider()
	if err != nil {
		return err
	}

	res, err := execer.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil && len(feedback) > 0 && feedback[0] {
		return err
	}

	if aff == 0 && len(feedback) > 0 && feedback[0] {
		return ErrZeroRowsAffected
	}

	return nil
}
