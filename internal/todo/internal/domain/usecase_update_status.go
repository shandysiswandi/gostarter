package domain

import (
	"context"
)

type UpdateStatus interface {
	Execute(ctx context.Context, in UpdateStatusInput) (*UpdateStatusOutput, error)
}

type UpdateStatusInput struct {
	ID     uint64 `validate:"gt=0"`
	Status string `validate:"required"`
}

type UpdateStatusOutput struct {
	ID     uint64
	Status TodoStatus
}
