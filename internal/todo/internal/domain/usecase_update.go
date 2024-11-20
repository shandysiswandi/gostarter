package domain

import (
	"context"
)

type Update interface {
	Call(ctx context.Context, in UpdateInput) (*Todo, error)
}

type UpdateInput struct {
	ID          uint64 `validate:"gt=0"`
	Title       string `validate:"min=5"`
	Description string `validate:"min=15"`
	Status      string `validate:"required"`
}
