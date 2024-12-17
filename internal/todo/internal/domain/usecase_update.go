package domain

import (
	"context"
)

type Update interface {
	Call(ctx context.Context, in UpdateInput) (*Todo, error)
}

type UpdateInput struct {
	ID          uint64 `validate:"required,gt=0"`
	Title       string `validate:"required,min=5"`
	Description string `validate:"required,min=15"`
	Status      string `validate:"required"`
}
