package domain

import (
	"context"
)

type Fetch interface {
	Execute(ctx context.Context, in FetchInput) ([]Todo, error)
}

type FetchInput struct {
	ID          string
	Title       string
	Description string
	Status      string
}
