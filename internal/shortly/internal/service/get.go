package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type GetStore interface {
	Get(ctx context.Context, key string) (*domain.Short, error)
}

type Get struct {
	store     GetStore
	validator validation.Validator
	logger    logger.Logger
}

func NewGet(store GetStore, v validation.Validator, l logger.Logger) *Get {
	return &Get{store: store, validator: v, logger: l}
}

func (g *Get) Call(ctx context.Context, in domain.GetInput) (*domain.GetOutput, error) {
	err := g.validator.Validate(in)
	if err != nil {
		g.logger.Error(ctx, "validation failed", err)

		return nil, err
	}

	resp, err := g.store.Get(ctx, in.Key)
	if err != nil {
		g.logger.Error(ctx, "failed to get", err)

		return nil, err
	}

	if resp == nil {
		g.logger.Warn(ctx, "data not found")

		return &domain.GetOutput{URL: ""}, nil
	}

	return &domain.GetOutput{URL: resp.URL}, nil
}
