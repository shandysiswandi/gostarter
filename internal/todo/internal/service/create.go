package service

import (
	"context"
	"errors"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/errs"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type CreateStore interface {
	Create(ctx context.Context, in entity.Todo) error
}

type Create struct {
	store     CreateStore
	uidnumber uid.NumberID
	validator validation.Validator
}

func NewCreate(store CreateStore, validator validation.Validator, uidnumber uid.NumberID) *Create {
	return &Create{
		store:     store,
		uidnumber: uidnumber,
		validator: validator,
	}
}

func (s *Create) Execute(ctx context.Context, in usecase.CreateInput) (*usecase.CreateOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		return nil, errs.WrapValidation("validation input fail", err)
	}

	id := s.uidnumber.Generate()

	err := s.store.Create(ctx, entity.Todo{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatusInitiate,
	})
	if errors.Is(err, entity.ErrTodoNotCreated) {
		return nil, errs.NewBusiness("failed to create todo")
	}

	if err != nil {
		return nil, errs.NewServerFrom(err)
	}

	return &usecase.CreateOutput{
		ID: id,
	}, nil
}
