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
	store     DeleteStore
	validator validation.Validator
}

func NewDelete(t *telemetry.Telemetry, s DeleteStore, v validation.Validator) *Delete {
	return &Delete{
		telemetry: t,
		store:     s,
		validator: v,
	}
}

func (s *Delete) Execute(ctx context.Context, in domain.DeleteInput) (*domain.DeleteOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	err := s.store.Delete(ctx, in.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to delete", err)

		return nil, goerror.NewServer("failed to delete todo", err)
	}

	return &domain.DeleteOutput{ID: in.ID}, nil
}
