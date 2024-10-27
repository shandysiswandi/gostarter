package domain

import "context"

type ResetPassword interface {
	Call(ctx context.Context, in ResetPasswordInput) (*ResetPasswordOutput, error)
}

type ResetPasswordInput struct {
	Token    string `validate:"required"`
	Password string `validate:"required"`
}

type ResetPasswordOutput struct{}
