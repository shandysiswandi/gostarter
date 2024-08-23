package usecase

import "context"

type Create interface {
	Execute(ctx context.Context, in CreateInput) (*CreateOutput, error)
}

type CreateInput struct {
	Title       string `validate:"min=5"`
	Description string `validate:"min=15"`
}

type CreateOutput struct {
	ID uint64
}
