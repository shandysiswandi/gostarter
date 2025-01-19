package domain

import "context"

type UpdatePassword interface {
	Call(ctx context.Context, in UpdatePasswordInput) (*User, error)
}

type UpdatePasswordInput struct {
	CurrentPassword string `validate:"required,min=5,max=60"`
	NewPassword     string `validate:"required,min=5,max=60"`
}
