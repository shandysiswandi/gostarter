package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type DeleteStore interface {
	Delete(ctx context.Context, in uint64) error
}

type Delete struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	store     DeleteStore
}

func NewDelete(dep Dependency, s DeleteStore) *Delete {
	return &Delete{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (s *Delete) Call(ctx context.Context, in domain.DeleteInput) (*domain.DeleteOutput, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.Delete")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	if err := s.store.Delete(ctx, in.ID); err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to delete", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.DeleteOutput{ID: in.ID}, nil
}
