package domain

import (
	"context"
	"time"
)

type Verify interface {
	Call(ctx context.Context, in VerifyInput) (*VerifyOutput, error)
}

type VerifyInput struct {
	Email string `validate:"required,email,min=5,max=100"`
	Code  string `validate:"required,length=6"`
}

type VerifyOutput struct {
	Email    string
	VerifyAt time.Time
}
