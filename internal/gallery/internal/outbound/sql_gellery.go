package outbound

import (
	"context"
	"database/sql"

	"github.com/shandysiswandi/gostarter/internal/gallery/internal/domain"
)

type SQLGallery struct {
	db *sql.DB
}

func NewSQLGallery(db *sql.DB) *SQLGallery {
	return &SQLGallery{db: db}
}

func (sg *SQLGallery) Upload(ctx context.Context, in domain.Image) error {
	return nil
}

func (sg *SQLGallery) GetImage(ctx context.Context, id uint64) (*domain.Image, error) {
	return nil, nil
}
