package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
)

type GetByID interface {
	Execute(ctx context.Context, in GetByIDInput) (*GetByIDOutput, error)
}

type GetByIDInput struct {
	ID uint64 `validate:"required,gte=0"`
}

type GetByIDOutput struct {
	ID          uint64
	Title       string
	Description string
	Status      entity.TodoStatus
}
