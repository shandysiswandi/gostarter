package usecase

import "context"

type Create interface {
	Execute(ctx context.Context, in CreateInput) (*CreateOutput, error)
}

type CreateInput struct {
	Title       string
	Description string
}

type CreateOutput struct {
	ID uint64
}
