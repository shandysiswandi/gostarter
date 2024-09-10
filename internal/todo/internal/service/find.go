package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type FindStore interface {
	Find(ctx context.Context, id uint64) (*domain.Todo, error)
}

type Find struct {
	log       logger.Logger
	store     FindStore
	validator validation.Validator
}

func NewFind(l logger.Logger, s FindStore, v validation.Validator) *Find {
	return &Find{
		log:       l,
		store:     s,
		validator: v,
	}
}

func (s *Find) Execute(ctx context.Context, in domain.FindInput) (*domain.Todo, error) {
	if err := s.validator.Validate(in); err != nil {
		s.log.Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	todo, err := s.store.Find(ctx, in.ID)
	if err != nil {
		s.log.Error(ctx, "todo fail to find", err)

		return nil, goerror.NewServer("failed to find todo", err)
	}

	if todo == nil {
		s.log.Warn(ctx, "todo is not found")

		return nil, goerror.NewBusiness("todo not found", goerror.CodeNotFound)
	}

	return &domain.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}, nil
}
