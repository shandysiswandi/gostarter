package domain

import (
	"context"
)

type DetachUserRole interface {
	Call(ctx context.Context, in DetachUserRoleInput) (*DetachUserRoleOutput, error)
}

type DetachUserRoleInput struct {
	UserID uint64 `validate:"required,gt=0"`
	RoleID uint64 `validate:"required,gt=0"`
}

type DetachUserRoleOutput struct {
	Message string
}
