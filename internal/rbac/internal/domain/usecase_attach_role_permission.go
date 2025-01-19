package domain

import (
	"context"
)

type AttachRolePermission interface {
	Call(ctx context.Context, in AttachRolePermissionInput) (*AttachRolePermissionOutput, error)
}

type AttachRolePermissionInput struct {
	RoleID        uint64   `validate:"required,gt=0"`
	PermissionIDs []uint64 `validate:"required,dive,gt=0"`
}

type AttachRolePermissionOutput struct {
	Message string
}
