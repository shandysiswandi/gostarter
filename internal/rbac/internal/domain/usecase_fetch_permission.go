package domain

import (
	"context"
)

type FetchPermission interface {
	Call(ctx context.Context, in FetchPermissionInput) (*FetchPermissionOutput, error)
}

type FetchPermissionInput struct {
	Cursor string
	Limit  string
	Name   string
}

type FetchPermissionOutput struct {
	Permissions []Permission
	NextCursor  string
	HasMore     bool
}
