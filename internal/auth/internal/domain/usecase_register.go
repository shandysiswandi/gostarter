package domain

import "context"

type Register interface {
	Call(ctx context.Context, in RegisterInput) (*RegisterOutput, error)
}

type RegisterInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterOutput struct{}
