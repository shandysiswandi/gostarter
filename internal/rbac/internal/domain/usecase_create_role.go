package domain

import (
	"context"
)

type CreateRole interface {
	Call(ctx context.Context, in CreateRoleInput) (*CreateRoleOutput, error)
}

type CreateRoleInput struct {
	Name        string `validate:"required,min=5,max=50"`
	Description string `validate:"required,min=15,max=255"`
}

type CreateRoleOutput struct {
	ID uint64
}
