package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
)

type FetchStore interface {
	Fetch(ctx context.Context, filter map[string]string) ([]domain.Todo, error)
}

type Fetch struct {
	log   logger.Logger
	store FetchStore
}

func NewFetch(l logger.Logger, s FetchStore) *Fetch {
	return &Fetch{
		log:   l,
		store: s,
	}
}

func (s *Fetch) Execute(ctx context.Context, in domain.FetchInput) ([]domain.Todo, error) {
	filter := map[string]string{
		"id":          in.ID,
		"title":       in.Title,
		"description": in.Description,
		"status":      in.Status,
	}

	todos, err := s.store.Fetch(ctx, filter)
	if err != nil {
		s.log.Error(ctx, "todo fail to fetch", err)

		return nil, goerror.NewServer("failed to fetch todo", err)
	}

	return todos, nil
}
