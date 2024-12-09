package domain

import "context"

type Logout interface {
	Call(ctx context.Context, in LogoutInput) (*LogoutOutput, error)
}

type LogoutInput struct {
	AccessToken string `validate:"required"`
}
type LogoutOutput struct {
	Message string
}
