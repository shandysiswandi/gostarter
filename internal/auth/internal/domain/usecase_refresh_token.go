package domain

import "context"

type RefreshToken interface {
	Call(ctx context.Context, in RefreshTokenInput) (*RefreshTokenOutput, error)
}

type RefreshTokenInput struct {
	RefreshToken string `validate:"required"`
}

type RefreshTokenOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64 // in seconds
	RefreshExpiresIn int64 // in seconds
}
