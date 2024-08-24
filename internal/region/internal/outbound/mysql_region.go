package outbound

import (
	"context"
	"database/sql"

	"github.com/shandysiswandi/gostarter/internal/region/internal/entity"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type MysqlRegion struct {
	db *sql.DB
}

func NewMysqlRegion(db *sql.DB) *MysqlRegion {
	return &MysqlRegion{
		db: db,
	}
}

func (m *MysqlRegion) Provinces(ctx context.Context, ids ...string) ([]entity.Province, error) {
	query := func() (string, []any, error) {
		ops := dbops.New().Select("id,name").From("provinces")
		if len(ids) > 0 {
			ops.WhereIn("id", ids...)
		}
		q, args := ops.ToSQL()

		return q, args, nil
	}

	return dbops.SQLGets[entity.Province](ctx, m.db, query)
}

func (m *MysqlRegion) Cities(ctx context.Context, pID string, ids ...string) ([]entity.City, error) {
	query := func() (string, []any, error) {
		ops := dbops.New().Select("id,province_id,name").From("cities")
		if len(ids) > 0 {
			ops.WhereIn("id", ids...)
		}

		if pID != "" {
			ops.Where("province_id", pID)
		}
		q, args := ops.ToSQL()

		return q, args, nil
	}

	return dbops.SQLGets[entity.City](ctx, m.db, query)
}

func (m *MysqlRegion) Districts(ctx context.Context, pID string, ids ...string) ([]entity.District, error) {
	query := func() (string, []any, error) {
		ops := dbops.New().Select("id,city_id,name").From("districts")
		if len(ids) > 0 {
			ops.WhereIn("id", ids...)
		}

		if pID != "" {
			ops.Where("city_id", pID)
		}
		q, args := ops.ToSQL()

		return q, args, nil
	}

	return dbops.SQLGets[entity.District](ctx, m.db, query)
}

func (m *MysqlRegion) Villages(ctx context.Context, pID string, ids ...string) ([]entity.Village, error) {
	query := func() (string, []any, error) {
		ops := dbops.New().Select("id,district_id,name").From("villages")
		if len(ids) > 0 {
			ops.WhereIn("id", ids...)
		}

		if pID != "" {
			ops.Where("district_id", pID)
		}
		q, args := ops.ToSQL()

		return q, args, nil
	}

	return dbops.SQLGets[entity.Village](ctx, m.db, query)
}
