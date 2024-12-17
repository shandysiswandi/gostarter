package domain

import "context"

type Delete interface {
	Call(ctx context.Context, in DeleteInput) (*DeleteOutput, error)
}

type DeleteInput struct {
	ID uint64 `validate:"required,gt=0"`
}

type DeleteOutput struct {
	ID uint64
}
