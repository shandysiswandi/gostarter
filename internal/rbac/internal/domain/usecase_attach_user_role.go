package domain

import (
	"context"
)

type AttachUserRole interface {
	Call(ctx context.Context, in AttachUserRoleInput) (*AttachUserRoleOutput, error)
}

type AttachUserRoleInput struct {
	UserID  uint64   `validate:"required,gt=0"`
	RoleIDs []uint64 `validate:"required,dive,gt=0"`
}

type AttachUserRoleOutput struct {
	Message string
}
