package domain

import "context"

type Create interface {
	Call(ctx context.Context, in CreateInput) (*CreateOutput, error)
}

type CreateInput struct {
	Title       string `validate:"required,min=5"`
	Description string `validate:"required,min=15"`
}

type CreateOutput struct {
	ID uint64
}
