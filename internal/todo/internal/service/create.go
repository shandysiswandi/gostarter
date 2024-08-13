package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type CreateStore interface {
	Create(ctx context.Context, in entity.Todo) error
}

type Create struct {
	store     CreateStore
	uidnumber uid.Number
	validator validation.Validator
}

func NewCreate(store CreateStore, validator validation.Validator, uidnumber uid.Number) *Create {
	return &Create{
		store:     store,
		uidnumber: uidnumber,
		validator: validator,
	}
}

func (s *Create) Execute(ctx context.Context, in usecase.CreateInput) (*usecase.CreateOutput, error) {
	id := s.uidnumber.Generate()
	err := s.store.Create(ctx, entity.Todo{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatusInitiate,
	})
	if err != nil {
		return nil, err
	}

	return &usecase.CreateOutput{
		ID: id,
	}, nil
}
