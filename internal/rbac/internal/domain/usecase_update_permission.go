package domain

import (
	"context"
)

type UpdatePermission interface {
	Call(ctx context.Context, in UpdatePermissionInput) (*UpdatePermissionOutput, error)
}

type UpdatePermissionInput struct {
	ID          uint64 `validate:"required,gt=0"`
	Name        string `validate:"required,min=5,max=50"`
	Description string `validate:"required,min=15,max=255"`
}

type UpdatePermissionOutput struct {
	ID          uint64
	Name        string
	Description string
}
