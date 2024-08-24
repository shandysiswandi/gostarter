package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/region/internal/entity"
)

type Search interface {
	Execute(ctx context.Context, in SearchInput) (*SearchOutput, error)
}

type SearchInput struct {
	By       string `validate:"omitempty,oneof=provinces cities districts villages"`
	ParentID string `validate:"omitempty,numeric"`
	IDs      string
}

type SearchOutput struct {
	Regions []entity.Region
}
