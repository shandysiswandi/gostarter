package domain

import (
	"context"
)

type Fetch interface {
	Call(ctx context.Context, in FetchInput) ([]Todo, error)
}

type FetchInput struct {
	ID          string
	Title       string
	Description string
	Status      string
}
