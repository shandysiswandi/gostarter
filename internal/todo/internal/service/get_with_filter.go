package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type GetWithFilterStore interface {
	GetWithFilter(ctx context.Context, filter map[string]string) ([]entity.Todo, error)
}

type GetWithFilter struct {
	store     GetWithFilterStore
	validator validation.Validator
}

func NewGetWithFilter(store GetWithFilterStore, validator validation.Validator) *GetWithFilter {
	return &GetWithFilter{
		store:     store,
		validator: validator,
	}
}

func (s *GetWithFilter) Execute(ctx context.Context, in usecase.GetWithFilterInput) (
	*usecase.GetWithFilterOutput, error,
) {
	filter := map[string]string{
		"id":          in.ID,
		"title":       in.Title,
		"description": in.Description,
		"status":      in.Status,
	}

	todos, err := s.store.GetWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &usecase.GetWithFilterOutput{
		Todos: todos,
	}, nil
}
