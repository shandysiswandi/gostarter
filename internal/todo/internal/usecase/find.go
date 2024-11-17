package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type FindStore interface {
	Find(ctx context.Context, id uint64) (*domain.Todo, error)
}

type Find struct {
	telemetry *telemetry.Telemetry
	store     FindStore
	validator validation.Validator
}

func NewFind(t *telemetry.Telemetry, s FindStore, v validation.Validator) *Find {
	return &Find{
		telemetry: t,
		store:     s,
		validator: v,
	}
}

func (s *Find) Execute(ctx context.Context, in domain.FindInput) (*domain.Todo, error) {
	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	todo, err := s.store.Find(ctx, in.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to find", err)

		return nil, goerror.NewServer("failed to find todo", err)
	}

	if todo == nil {
		s.telemetry.Logger().Warn(ctx, "todo is not found")

		return nil, goerror.NewBusiness("todo not found", goerror.CodeNotFound)
	}

	return &domain.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}, nil
}
