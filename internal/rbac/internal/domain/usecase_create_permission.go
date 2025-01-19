package domain

import (
	"context"
)

type CreatePermission interface {
	Call(ctx context.Context, in CreatePermissionInput) (*CreatePermissionOutput, error)
}

type CreatePermissionInput struct {
	Name        string `validate:"required,min=5,max=50"`
	Description string `validate:"required,min=15,max=255"`
}

type CreatePermissionOutput struct {
	ID uint64
}
