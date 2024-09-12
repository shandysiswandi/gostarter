package domain

import "context"

type GetImage interface {
	Call(ctx context.Context, in GetImageInput) (*GetImageOutput, error)
}

type GetImageInput struct {
	ID uint64 `validate:"required"`
}

type GetImageOutput struct {
	URLImage string
}
