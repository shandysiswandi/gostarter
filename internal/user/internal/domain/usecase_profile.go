package domain

import "context"

type Profile interface {
	Call(ctx context.Context, in ProfileInput) (*User, error)
}

type ProfileInput struct {
	Email string `validate:"required,email"`
}
