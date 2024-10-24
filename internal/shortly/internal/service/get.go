package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type GetStore interface {
	Get(ctx context.Context, key string) (*domain.Short, error)
}

type Get struct {
	store     GetStore
	validator validation.Validator
	telemetry *telemetry.Telemetry
}

func NewGet(store GetStore, v validation.Validator, t *telemetry.Telemetry) *Get {
	return &Get{store: store, validator: v, telemetry: t}
}

func (g *Get) Call(ctx context.Context, in domain.GetInput) (*domain.GetOutput, error) {
	err := g.validator.Validate(in)
	if err != nil {
		g.telemetry.Logger().Error(ctx, "validation failed", err)

		return nil, err
	}

	resp, err := g.store.Get(ctx, in.Key)
	if err != nil {
		g.telemetry.Logger().Error(ctx, "failed to get", err)

		return nil, err
	}

	if resp == nil {
		g.telemetry.Logger().Warn(ctx, "data not found")

		return &domain.GetOutput{URL: ""}, nil
	}

	return &domain.GetOutput{URL: resp.URL}, nil
}
