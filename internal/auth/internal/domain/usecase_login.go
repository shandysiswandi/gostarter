package domain

import "context"

type Login interface {
	Call(ctx context.Context, in LoginInput) (*LoginOutput, error)
}

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64 // in seconds
	RefreshExpiresIn int64 // in seconds
}
