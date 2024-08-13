package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
)

type UpdateStatus interface {
	Execute(ctx context.Context, in UpdateStatusInput) (*UpdateStatusOutput, error)
}

type UpdateStatusInput struct {
	ID     uint64
	Status string
}

type UpdateStatusOutput struct {
	ID     uint64
	Status entity.TodoStatus
}
