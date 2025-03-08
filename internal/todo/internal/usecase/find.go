package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type FindStore interface {
	Find(ctx context.Context, id uint64) (*domain.Todo, error)
}

type Find struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	store     FindStore
}

func NewFind(dep Dependency, s FindStore) *Find {
	return &Find{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (s *Find) Call(ctx context.Context, in domain.FindInput) (*domain.Todo, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.Find")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	todo, err := s.store.Find(ctx, in.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to find", err)

		return nil, goerror.NewServerInternal(err)
	}

	if todo == nil {
		s.telemetry.Logger().Warn(ctx, "todo is not found")

		return nil, goerror.NewBusiness("todo not found", goerror.CodeNotFound)
	}

	return &domain.Todo{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}, nil
}
