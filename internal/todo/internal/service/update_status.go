package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	uc "github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/errs"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStatusStore interface {
	UpdateStatus(ctx context.Context, in uint64, status entity.TodoStatus) error
}

type UpdateStatus struct {
	store     UpdateStatusStore
	validator validation.Validator
}

func NewUpdateStatus(store UpdateStatusStore, validator validation.Validator) *UpdateStatus {
	return &UpdateStatus{
		store:     store,
		validator: validator,
	}
}

func (s *UpdateStatus) Execute(ctx context.Context, in uc.UpdateStatusInput) (*uc.UpdateStatusOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		return nil, errs.WrapValidation("validation input fail", err)
	}

	sts := entity.ParseTodoStatus(in.Status)

	err := s.store.UpdateStatus(ctx, in.ID, sts)
	if err != nil {
		return nil, errs.NewServerFrom(err)
	}

	return &uc.UpdateStatusOutput{
		ID:     in.ID,
		Status: sts,
	}, nil
}
