package domain

import (
	"context"
)

type Search interface {
	Call(ctx context.Context, in SearchInput) ([]Region, error)
}

type SearchInput struct {
	By       string `validate:"omitempty,oneof=provinces cities districts villages"`
	ParentID string `validate:"omitempty,numeric"`
	IDs      string
}
