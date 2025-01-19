package domain

import (
	"context"
)

type UpdateRole interface {
	Call(ctx context.Context, in UpdateRoleInput) (*UpdateRoleOutput, error)
}

type UpdateRoleInput struct {
	ID          uint64 `validate:"required,gt=0"`
	Name        string `validate:"required,min=5,max=50"`
	Description string `validate:"required,min=15,max=255"`
}

type UpdateRoleOutput struct {
	ID          uint64
	Name        string
	Description string
}
