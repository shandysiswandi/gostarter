package domain

import "context"

type FindRole interface {
	Call(ctx context.Context, in FindRoleInput) (*Role, error)
}

type FindRoleInput struct {
	ID uint64 `validate:"required,gt=0"`
}
