package service

import (
	"context"
	"errors"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type CreateStore interface {
	Create(ctx context.Context, in domain.Todo) error
}

type Create struct {
	telemetry *telemetry.Telemetry
	store     CreateStore
	uidnumber uid.NumberID
	validator validation.Validator
}

func NewCreate(t *telemetry.Telemetry, s CreateStore, v validation.Validator, idgen uid.NumberID) *Create {
	return &Create{
		telemetry: t,
		store:     s,
		uidnumber: idgen,
		validator: v,
	}
}

func (s *Create) Execute(ctx context.Context, in domain.CreateInput) (*domain.CreateOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	id := s.uidnumber.Generate()

	err := s.store.Create(ctx, domain.Todo{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      domain.TodoStatusInitiate,
	})
	if errors.Is(err, domain.ErrTodoNotCreated) {
		s.telemetry.Logger().Warn(ctx, "todo created but db not affected")

		return nil, goerror.NewBusiness("failed to create todo", goerror.CodeUnknown)
	}

	if err != nil {
		s.telemetry.Logger().Error(ctx, "todo fail to create", err)

		return nil, goerror.NewServer("failed to create todo", err)
	}

	return &domain.CreateOutput{ID: id}, nil
}
