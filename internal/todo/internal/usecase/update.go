package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
)

type Update interface {
	Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error)
}

type UpdateInput struct {
	ID          uint64 `validate:"gt=0"`
	Title       string `validate:"min=5"`
	Description string `validate:"min=15"`
	Status      string `validate:"required"`
}

type UpdateOutput struct {
	ID          uint64
	Title       string
	Description string
	Status      entity.TodoStatus
}
