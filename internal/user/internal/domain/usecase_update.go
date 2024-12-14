package domain

import "context"

type Update interface {
	Call(ctx context.Context, in UpdateInput) (*User, error)
}

type UpdateInput struct {
	ID   uint64 `validate:"required"`
	Name string `validate:"required,min=5,max=100"`
}
