package domain

import "context"

type ForgotPassword interface {
	Call(ctx context.Context, in ForgotPasswordInput) (*ForgotPasswordOutput, error)
}

type ForgotPasswordInput struct {
	Email string `validate:"required,email,min=5,max=100"`
}

type ForgotPasswordOutput struct {
	Email   string
	Message string
}
