package sqlkit

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type Execer interface {
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

type Result struct {
	LastInsertID uint64
	RowsAffected uint64
}

func Update[M Model](ctx context.Context, db *DB, data map[string]any, exps ...Expression) (Result, error) {
	var m M

	return resultOrError(db.qb.Update(m.Table()).Set(data).Where(exps...).Executor().ExecContext(ctx))
}

func Delete[M Model](ctx context.Context, db *DB, exps ...Expression) (Result, error) {
	var m M

	return resultOrError(db.qb.Delete(m.Table()).Where(exps...).Executor().ExecContext(ctx))
}

func Exec(ctx context.Context, db *DB, query string, args ...any) (Result, error) {
	return resultOrError(db.Execer(ctx).ExecContext(ctx, query, args...))
}

func resultOrError(result sql.Result, err error) (Result, error) {
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

	return Result{LastInsertID: uint64(lastId), RowsAffected: uint64(aff)}, nil
}

func AnyToValue(args []any) []driver.Value {
	driverArgs := make([]driver.Value, 0, len(args))
	for _, arg := range args {
		driverArgs = append(driverArgs, arg)
	}

	return driverArgs
}
