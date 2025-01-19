package domain

import (
	"context"
)

type FetchRole interface {
	Call(ctx context.Context, in FetchRoleInput) (*FetchRoleOutput, error)
}

type FetchRoleInput struct {
	Cursor string
	Limit  string
	Name   string
}

type FetchRoleOutput struct {
	Roles      []Role
	NextCursor string
	HasMore    bool
}
