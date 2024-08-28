package service

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type SetStore interface {
	Set(ctx context.Context, value domain.Short) error
}

type Set struct {
	store     SetStore
	validator validation.Validator
	logger    logger.Logger
}

func NewSet(store SetStore, v validation.Validator, l logger.Logger) *Set {
	return &Set{store: store, validator: v, logger: l}
}

func (g *Set) Call(ctx context.Context, in domain.SetInput) (*domain.SetOutput, error) {
	now := time.Now()
	theKey := base64.URLEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(now.UnixNano()), 10)))

	err := g.validator.Validate(in)
	if err != nil {
		g.logger.Error(ctx, "validation failed", err)

		return nil, err
	}

	err = g.store.Set(ctx, domain.Short{
		Key:     theKey,
		URL:     in.URL,
		Slug:    in.Slug,
		Expired: now,
	})
	if err != nil {
		g.logger.Error(ctx, "failed to save", err)

		return nil, err
	}

	return &domain.SetOutput{Key: theKey}, nil
}
