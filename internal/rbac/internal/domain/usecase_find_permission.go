package domain

import (
	"context"
)

type FindPermission interface {
	Call(ctx context.Context, in FindPermissionInput) (*Permission, error)
}

type FindPermissionInput struct {
	ID uint64 `validate:"required,gt=0"`
}
