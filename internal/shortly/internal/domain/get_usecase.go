package domain

import "context"

type Get interface {
	Call(ctx context.Context, in GetInput) (*GetOutput, error)
}

type GetInput struct {
	Key string `validate:"required"`
}

type GetOutput struct {
	URL string
}
