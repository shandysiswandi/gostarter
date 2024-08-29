package domain

import "context"

type Set interface {
	Call(ctx context.Context, in SetInput) (*SetOutput, error)
}

type SetInput struct {
	URL  string `validate:"required,url"`
	Slug string
}

type SetOutput struct {
	Key string
}
