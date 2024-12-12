package domain

import (
	"context"
)

type Fetch interface {
	Call(ctx context.Context, in FetchInput) (*FetchOutput, error)
}

type FetchInput struct {
	Cursor string
	Limit  string
	Status string
}

type FetchOutput struct {
	Todos      []Todo
	NextCursor string
	HasMore    bool
}
