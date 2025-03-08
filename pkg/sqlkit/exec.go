package sqlkit

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/shandysiswandi/goreng/telemetry/logger"
)

type Execer interface {
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

type Result struct {
	LastInsertID uint64
	RowsAffected uint64
}

func Exec(ctx context.Context, db *DB, query string, args ...any) (Result, error) {
	result, err := db.Execer(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return Result{}, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Result{}, err
	}

	aff, err := result.RowsAffected()
	if err != nil {
		return Result{LastInsertID: uint64(lastId)}, err
	}

	db.log.Info(ctx, "Exec Result", logger.KeyVal("result", Result{
		LastInsertID: uint64(lastId),
		RowsAffected: uint64(aff),
	}))

	return Result{LastInsertID: uint64(lastId), RowsAffected: uint64(aff)}, nil
}

func AnyToValue(args []any) []driver.Value {
	driverArgs := make([]driver.Value, 0, len(args))
	for _, arg := range args {
		driverArgs = append(driverArgs, arg)
	}

	return driverArgs
}
