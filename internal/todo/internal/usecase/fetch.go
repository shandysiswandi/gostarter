package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type FetchStore interface {
	Fetch(ctx context.Context, filter map[string]string) ([]domain.Todo, error)
}

type Fetch struct {
	telemetry *telemetry.Telemetry
	store     FetchStore
}

func NewFetch(dep Dependency, s FetchStore) *Fetch {
	return &Fetch{
		telemetry: dep.Telemetry,
		store:     s,
	}
}

func (s *Fetch) Call(ctx context.Context, in domain.FetchInput) ([]domain.Todo, error) {
	filter := map[string]string{
		"id":          in.ID,
		"title":       in.Title,
		"description": in.Description,
		"status":      in.Status,
	}

	todos, err := s.store.Fetch(ctx, filter)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to fetch", err)

		return nil, goerror.NewServer("failed to fetch todo", err)
	}

	return todos, nil
}
