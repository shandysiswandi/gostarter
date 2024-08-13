package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStore interface {
	Update(ctx context.Context, in entity.Todo) error
}

type Update struct {
	store     UpdateStore
	validator validation.Validator
}

func NewUpdate(store UpdateStore, validator validation.Validator) *Update {
	return &Update{
		store:     store,
		validator: validator,
	}
}

func (s *Update) Execute(ctx context.Context, in usecase.UpdateInput) (*usecase.UpdateOutput, error) {
	sts := entity.ParseTodoStatus(in.Status)
	err := s.store.Update(ctx, entity.Todo{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	})
	if err != nil {
		return nil, err
	}

	return &usecase.UpdateOutput{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	}, nil
}
