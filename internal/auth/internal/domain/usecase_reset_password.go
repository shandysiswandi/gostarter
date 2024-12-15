package domain

import "context"

type ResetPassword interface {
	Call(ctx context.Context, in ResetPasswordInput) (*ResetPasswordOutput, error)
}

type ResetPasswordInput struct {
	Token    string `validate:"required,min=5"`
	Password string `validate:"required,min=5,max=60"`
}

type ResetPasswordOutput struct {
	Message string
}
