package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
)

type ProfileStore interface {
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Profile struct {
	tel   *telemetry.Telemetry
	store ProfileStore
}

func NewProfile(dep Dependency, s ProfileStore) *Profile {
	return &Profile{
		tel:   dep.Telemetry,
		store: s,
	}
}

func (p *Profile) Call(ctx context.Context, _ domain.ProfileInput) (*domain.User, error) {
	ctx, span := p.tel.Tracer().Start(ctx, "user.usecase.Profile")
	defer span.End()

	var email string
	if clm := lib.GetJWTClaim(ctx); clm != nil {
		email = clm.Subject
	}

	user, err := p.store.UserByEmail(ctx, email)
	if err != nil {
		p.tel.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", email))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		p.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("email", email))

		return nil, goerror.NewBusiness("user not found", goerror.CodeNotFound)
	}
	user.Password = "***" // re-assign or mask the password for security reason

	return &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
