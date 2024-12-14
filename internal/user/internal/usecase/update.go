package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdateStore interface {
	Update(ctx context.Context, user map[string]any) error
}

type Update struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	store     UpdateStore
}

func NewUpdate(dep Dependency, s UpdateStore) *Update {
	return &Update{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (p *Update) Call(ctx context.Context, in domain.UpdateInput) (*domain.User, error) {
	ctx, span := p.tel.Tracer().Start(ctx, "usecase.Update")
	defer span.End()

	if err := p.validator.Validate(in); err != nil {
		p.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	var email string
	if clm := jwt.GetClaim(ctx); clm != nil {
		email = clm.Subject
	}

	user := map[string]any{"id": in.ID, "name": in.Name}
	if err := p.store.Update(ctx, user); err != nil {
		p.tel.Logger().Error(ctx, "failed to update user", err, logger.KeyVal("id", in.ID))

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.User{
		ID:       in.ID,
		Name:     in.Name,
		Email:    email,
		Password: "***",
	}, nil
}
