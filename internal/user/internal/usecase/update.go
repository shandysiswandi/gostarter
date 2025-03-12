package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
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
	ctx, span := p.tel.Tracer().Start(ctx, "user.usecase.Update")
	defer span.End()

	if err := p.validator.Validate(in); err != nil {
		p.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	var email string
	var uid uint64
	if clm := lib.GetJWTClaim(ctx); clm != nil {
		email = clm.Subject
		uid = clm.AuthID
	}

	user := map[string]any{"id": uid, "name": in.Name}
	if err := p.store.Update(ctx, user); err != nil {
		p.tel.Logger().Error(ctx, "failed to update user", err, logger.KeyVal("id", uid))

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.User{
		ID:       uid,
		Name:     in.Name,
		Email:    email,
		Password: "***",
	}, nil
}
