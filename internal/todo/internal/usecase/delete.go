package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
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

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	if err := s.store.Delete(ctx, in.ID); err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to delete", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.DeleteOutput{ID: in.ID}, nil
}
