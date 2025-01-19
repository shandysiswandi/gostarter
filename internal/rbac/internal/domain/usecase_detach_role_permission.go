package domain

import (
	"context"
)

type DetachRolePermission interface {
	Call(ctx context.Context, in DetachRolePermissionInput) (*DetachRolePermissionOutput, error)
}

type DetachRolePermissionInput struct {
	RoleID       uint64 `validate:"required,gt=0"`
	PermissionID uint64 `validate:"required,gt=0"`
}

type DetachRolePermissionOutput struct {
	Message string
}
