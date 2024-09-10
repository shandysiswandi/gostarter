package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type DeleteStore interface {
	Delete(ctx context.Context, in uint64) error
}

type Delete struct {
	log       logger.Logger
	store     DeleteStore
	validator validation.Validator
}

func NewDelete(l logger.Logger, s DeleteStore, v validation.Validator) *Delete {
	return &Delete{
		log:       l,
		store:     s,
		validator: v,
	}
}

func (s *Delete) Execute(ctx context.Context, in domain.DeleteInput) (*domain.DeleteOutput, error) {
	if err := s.validator.Validate(in); err != nil {
		s.log.Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	err := s.store.Delete(ctx, in.ID)
	if err != nil {
		s.log.Error(ctx, "todo fail to delete", err)

		return nil, goerror.NewServer("failed to delete todo", err)
	}

	return &domain.DeleteOutput{ID: in.ID}, nil
}
