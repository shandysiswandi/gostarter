package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
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

func (s *UpdateStatus) Execute(ctx context.Context, in usecase.UpdateStatusInput) (
	*usecase.UpdateStatusOutput, error,
) {
	sts := entity.ParseTodoStatus(in.Status)

	err := s.store.UpdateStatus(ctx, in.ID, sts)
	if err != nil {
		return nil, err
	}

	return &usecase.UpdateStatusOutput{
		ID:     in.ID,
		Status: sts,
	}, nil
}
