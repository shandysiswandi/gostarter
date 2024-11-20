package domain

import (
	"context"
)

type Find interface {
	Call(ctx context.Context, in FindInput) (*Todo, error)
}

type FindInput struct {
	ID uint64 `validate:"gt=0"`
}
