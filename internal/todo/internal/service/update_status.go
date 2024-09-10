package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStatusStore interface {
	UpdateStatus(ctx context.Context, in uint64, status domain.TodoStatus) error
}

type UpdateStatus struct {
	log       logger.Logger
	store     UpdateStatusStore
	validator validation.Validator
}

func NewUpdateStatus(l logger.Logger, s UpdateStatusStore, v validation.Validator) *UpdateStatus {
	return &UpdateStatus{
		log:       l,
		store:     s,
		validator: v,
	}
}

func (s *UpdateStatus) Execute(ctx context.Context, in domain.UpdateStatusInput) (*domain.UpdateStatusOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		s.log.Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	sts := domain.ParseTodoStatus(in.Status)

	err := s.store.UpdateStatus(ctx, in.ID, sts)
	if err != nil {
		s.log.Error(ctx, "todo fail to update status", err)

		return nil, goerror.NewServer("failed to update status todo", err)
	}

	return &domain.UpdateStatusOutput{ID: in.ID, Status: sts}, nil
}
