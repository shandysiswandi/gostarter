package domain

import (
	"context"

	"github.com/shandysiswandi/gostarter/pkg/enum"
)

type UpdateStatus interface {
	Call(ctx context.Context, in UpdateStatusInput) (*UpdateStatusOutput, error)
}

type UpdateStatusInput struct {
	ID     uint64 `validate:"gt=0"`
	Status string `validate:"required"`
}

type UpdateStatusOutput struct {
	ID     uint64
	Status enum.Enum[TodoStatus]
}
