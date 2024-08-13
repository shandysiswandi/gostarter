package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
)

type GetWithFilter interface {
	Execute(ctx context.Context, in GetWithFilterInput) (*GetWithFilterOutput, error)
}

type GetWithFilterInput struct {
	ID          string
	Title       string
	Description string
	Status      string
}

type GetWithFilterOutput struct {
	Todos []entity.Todo
}
