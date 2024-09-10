package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStore interface {
	Update(ctx context.Context, in domain.Todo) error
}

type Update struct {
	log       logger.Logger
	store     UpdateStore
	validator validation.Validator
}

func NewUpdate(l logger.Logger, s UpdateStore, v validation.Validator) *Update {
	return &Update{
		log:       l,
		store:     s,
		validator: v,
	}
}

func (s *Update) Execute(ctx context.Context, in domain.UpdateInput) (*domain.Todo, error) {
	if err := s.validator.Validate(in); err != nil {
		s.log.Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	sts := domain.ParseTodoStatus(in.Status)

	err := s.store.Update(ctx, domain.Todo{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	})
	if err != nil {
		s.log.Error(ctx, "todo fail to update", err)

		return nil, goerror.NewServer("failed to update todo", err)
	}

	return &domain.Todo{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Status:      sts,
	}, nil
}
