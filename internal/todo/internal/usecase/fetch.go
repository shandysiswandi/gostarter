package usecase

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/pagination"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type FetchStore interface {
	Fetch(ctx context.Context, filter map[string]any) ([]domain.Todo, error)
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

func (s *Fetch) Call(ctx context.Context, in domain.FetchInput) (*domain.FetchOutput, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.Fetch")
	defer span.End()

	cursor, limit := pagination.ParseCursorBased(in.Cursor, in.Limit)

	filter := map[string]any{
		"cursor": cursor,
		"limit":  limit,
	}

	if in.Status != "" {
		filter["status"] = enum.New(enum.Parse[domain.TodoStatus](in.Status))
	}

	todos, err := s.store.Fetch(ctx, filter)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to fetch", err)

		return nil, goerror.NewServerInternal(err)
	}

	nextCursor := ""
	hasMore := len(todos) > limit

	if hasMore {
		nextCursor = base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatUint(todos[limit].ID, 10)))
		todos = todos[:limit]
	}

	return &domain.FetchOutput{
		Todos:      todos,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
