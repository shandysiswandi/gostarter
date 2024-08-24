package outbound

import (
	"context"
	"database/sql"

	"github.com/shandysiswandi/gostarter/internal/region/internal/entity"
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
	return nil, nil
}

func (m *MysqlRegion) Cities(ctx context.Context, parentID string, ids ...string) ([]entity.City, error) {
	return nil, nil
}

func (m *MysqlRegion) Districts(ctx context.Context, parentID string, ids ...string) ([]entity.District, error) {
	return nil, nil
}

func (m *MysqlRegion) Villages(ctx context.Context, parentID string, ids ...string) ([]entity.Village, error) {
	return nil, nil
}
