package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type ProfileStore interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Profile struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	store     ProfileStore
}

func NewProfile(dep Dependency, s ProfileStore) *Profile {
	return &Profile{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		store:     s,
	}
}

func (p *Profile) Call(ctx context.Context, in domain.ProfileInput) (*domain.User, error) {
	ctx, span := p.tel.Tracer().Start(ctx, "usecase.Profile")
	defer span.End()

	if err := p.validator.Validate(in); err != nil {
		p.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	user, err := p.store.FindUserByEmail(ctx, in.Email)
	if err != nil {
		p.tel.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		p.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("user not found", goerror.CodeNotFound)
	}
	user.Password = "***" // re-assign or mask the password for security reason

	return &domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
