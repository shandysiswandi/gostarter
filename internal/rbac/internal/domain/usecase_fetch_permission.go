package domain

import (
	"context"
)

type FetchPermission interface {
	Call(ctx context.Context, in FetchPermissionInput) (*FetchPermissionOutput, error)
}

type FetchPermissionInput struct {
	Name string `validate:"max=50"`
}

type FetchPermissionOutput struct {
	Permissions []Permission
	NextCursor  string
	HasMore     bool
}
