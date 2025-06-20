package domain

import "context"

type Login interface {
	Call(ctx context.Context, in LoginInput) (*LoginOutput, error)
}

type LoginInput struct {
	Email    string `validate:"required,email,min=5,max=100"`
	Password string `validate:"required,min=8,max=60"`
}

type LoginOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64 // in seconds
	RefreshExpiresIn int64 // in seconds
}
