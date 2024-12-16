package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStatusStore interface {
	UpdateStatus(ctx context.Context, in uint64, sts enum.Enum[domain.TodoStatus]) error
}

type UpdateStatus struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	store     UpdateStatusStore
}

func NewUpdateStatus(dep Dependency, s UpdateStatusStore) *UpdateStatus {
	return &UpdateStatus{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (s *UpdateStatus) Call(ctx context.Context, in domain.UpdateStatusInput) (
	*domain.UpdateStatusOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.UpdateStatus")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	sts := enum.New(enum.Parse[domain.TodoStatus](in.Status))
	err := s.store.UpdateStatus(ctx, in.ID, sts)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to update status", err)

		return nil, goerror.NewServer("failed to update status todo", err)
	}

	return &domain.UpdateStatusOutput{
		ID:     in.ID,
		Status: sts,
	}, nil
}
