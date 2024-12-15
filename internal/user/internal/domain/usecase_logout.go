package domain

import "context"

type Logout interface {
	Call(ctx context.Context, in LogoutInput) (*LogoutOutput, error)
}

type LogoutInput struct {
	AccessToken string `validate:"required,min=5"`
}
type LogoutOutput struct {
	Message string
}
