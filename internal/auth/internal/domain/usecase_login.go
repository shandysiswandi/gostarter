package domain

import "context"

type Login interface {
	Call(ctx context.Context, in LoginInput) (LoginOutput, error)
}

type LoginInput struct{}

type LoginOutput struct{}
