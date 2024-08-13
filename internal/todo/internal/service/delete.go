package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type DeleteStore interface {
	Delete(ctx context.Context, in uint64) error
}

type Delete struct {
	store     DeleteStore
	validator validation.Validator
}

func NewDelete(store DeleteStore, validator validation.Validator) *Delete {
	return &Delete{
		store:     store,
		validator: validator,
	}
}

func (s *Delete) Execute(ctx context.Context, in usecase.DeleteInput) (*usecase.DeleteOutput, error) {
	err := s.store.Delete(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &usecase.DeleteOutput{
		ID: in.ID,
	}, nil
}
