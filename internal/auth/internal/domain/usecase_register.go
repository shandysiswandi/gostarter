package domain

import "context"

type Register interface {
	Call(ctx context.Context, in RegisterInput) (*RegisterOutput, error)
}

type RegisterInput struct {
	Name     string `validate:"required,min=5,max=100"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RegisterOutput struct {
	Email string
}
