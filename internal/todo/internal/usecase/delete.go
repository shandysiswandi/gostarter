package usecase

import "context"

type Delete interface {
	Execute(ctx context.Context, in DeleteInput) (*DeleteOutput, error)
}

type DeleteInput struct {
	ID uint64
}

type DeleteOutput struct {
	ID uint64
}
