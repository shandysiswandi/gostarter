package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type GetByIDStore interface {
	GetByID(ctx context.Context, id uint64) (*entity.Todo, error)
}

type GetByID struct {
	store     GetByIDStore
	validator validation.Validator
}

func NewGetByID(store GetByIDStore, validator validation.Validator) *GetByID {
	return &GetByID{
		store:     store,
		validator: validator,
	}
}

func (s *GetByID) Execute(ctx context.Context, in usecase.GetByIDInput) (*usecase.GetByIDOutput, error) {
	todo, err := s.store.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, entity.ErrTodoNotFound
	}

	return &usecase.GetByIDOutput{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}, nil
}
