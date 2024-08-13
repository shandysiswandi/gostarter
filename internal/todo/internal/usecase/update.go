package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
)

type Update interface {
	Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error)
}

type UpdateInput struct {
	ID          uint64
	Title       string
	Description string
	Status      string
}

type UpdateOutput struct {
	ID          uint64
	Title       string
	Description string
	Status      entity.TodoStatus
}
