package domain

import "context"

type ForgotPassword interface {
	Call(ctx context.Context, in ForgotPasswordInput) (*ForgotPasswordOutput, error)
}

type ForgotPasswordInput struct {
	Email string `validate:"required,email"`
}

type ForgotPasswordOutput struct{}
