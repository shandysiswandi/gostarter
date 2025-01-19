package domain

import (
	"context"
)

type FetchRole interface {
	Call(ctx context.Context, in FetchRoleInput) (*FetchRoleOutput, error)
}

type FetchRoleInput struct {
	Name string `validate:"max=50"`
}

type FetchRoleOutput struct {
	Roles      []Role
	NextCursor string
	HasMore    bool
}
