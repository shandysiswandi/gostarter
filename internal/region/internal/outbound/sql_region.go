package outbound

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shandysiswandi/gostarter/internal/region/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type SQLRegion struct {
	db     *sql.DB
	qu     goqu.DialectWrapper
	config config.Config
}

func NewSQLRegion(db *sql.DB, config config.Config) *SQLRegion {
	qu := goqu.Dialect("mysql")
	if config.GetString("database.driver") == "postgres" {
		qu = goqu.Dialect("postgres")
	}

	return &SQLRegion{
		db:     db,
		qu:     qu,
		config: config,
	}
}

func (s *SQLRegion) Provinces(ctx context.Context, ids ...string) ([]domain.Province, error) {
	query := func() (string, []any, error) {
		ops := s.qu.Select("id", "name").From("provinces").Limit(10)

		if len(ids) > 0 {
			ops = ops.Where(goqu.Ex{"id": ids})
		}

		return ops.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Province](ctx, s.db, query)
}

func (s *SQLRegion) Cities(ctx context.Context, pID string, ids ...string) ([]domain.City, error) {
	query := func() (string, []any, error) {
		ops := s.qu.Select("id", "name").From("cities").Limit(10)

		if len(ids) > 0 {
			ops = ops.Where(exp.Ex{"id": ids})
		}

		if pID != "" {
			ops = ops.Where(exp.Ex{"province_id": pID})
		}

		return ops.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.City](ctx, s.db, query)
}

func (s *SQLRegion) Districts(ctx context.Context, pID string, ids ...string) ([]domain.District, error) {
	query := func() (string, []any, error) {
		ops := s.qu.Select("id", "name").From("districts").Limit(10)

		if len(ids) > 0 {
			ops = ops.Where(exp.Ex{"id": ids})
		}

		if pID != "" {
			ops = ops.Where(exp.Ex{"city_id": pID})
		}

		return ops.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.District](ctx, s.db, query)
}

func (s *SQLRegion) Villages(ctx context.Context, pID string, ids ...string) ([]domain.Village, error) {
	query := func() (string, []any, error) {
		ops := s.qu.Select("id", "name").From("villages").Limit(10)

		if len(ids) > 0 {
			ops = ops.Where(exp.Ex{"id": ids})
		}

		if pID != "" {
			ops = ops.Where(exp.Ex{"district_id": pID})
		}

		return ops.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Village](ctx, s.db, query)
}
