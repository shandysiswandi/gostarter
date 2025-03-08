package sqlkit

import (
	"context"
	"database/sql"

	"github.com/shandysiswandi/goreng/telemetry/logger"
)

type Querier interface {
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
}

func One[M Model](ctx context.Context, db *DB, exps ...Expression) (*M, error) {
	var m M

	qq := db.qb.From(m.Table()).Select(m).Where(exps...).Limit(1)
	query, _, _ := qq.ToSQL()

	db.log.Debug(ctx, "get one data", logger.KeyVal("query", query))

	found, err := qq.ScanStructContext(ctx, &m)
	if err != nil {
		db.log.Error(ctx, "error when scan rows to destination", err)

		return nil, err
	}

	if !found {
		db.log.Warn(ctx, "get one data is not found", logger.KeyVal("query", query))

		return nil, nil //nolint:nilnil // just a nil result
	}

	return &m, nil
}

func Many[M Model](ctx context.Context, db *DB, exps ...Expression) ([]M, error) {
	var results []M
	var m M

	query, args, err := db.qb.From(m.Table()).Where(exps...).ToSQL()
	if err != nil {
		db.log.Error(ctx, "error when generate sql", err)

		return nil, err
	}

	if err := db.Scan(ctx, &results, query, args...); err != nil {
		db.log.Error(ctx, "error when scan rows to destination", err, logger.KeyVal("query", query))

		return nil, ErrScanRow
	}

	return results, nil
}
