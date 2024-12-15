package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStore interface {
	Update(ctx context.Context, in domain.Todo) error
}

type Update struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	store     UpdateStore
}

func NewUpdate(dep Dependency, s UpdateStore) *Update {
	return &Update{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (s *Update) Call(ctx context.Context, in domain.UpdateInput) (*domain.Todo, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "todo.usecase.Update")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	sts := domain.ParseTodoStatus(in.Status)
	userID := uint64(0)
	if clm := jwt.GetClaim(ctx); clm != nil {
		userID = clm.AuthID
	}

	err := s.store.Update(ctx, domain.Todo{
		ID:          in.ID,
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	})
	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to update", err)

		return nil, goerror.NewServer("failed to update todo", err)
	}

	return &domain.Todo{
		ID:          in.ID,
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	}, nil
}
