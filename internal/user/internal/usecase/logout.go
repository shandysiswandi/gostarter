package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type LogoutStore interface {
	DeleteTokenByAccess(ctx context.Context, token string) error
}

type Logout struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	secHash   hash.Hash
	store     LogoutStore
}

func NewLogout(dep Dependency, s LogoutStore) *Logout {
	return &Logout{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		secHash:   dep.SecHash,
		store:     s,
	}
}

func (l *Logout) Call(ctx context.Context, in domain.LogoutInput) (*domain.LogoutOutput, error) {
	ctx, span := l.tel.Tracer().Start(ctx, "usecase.Logout")
	defer span.End()

	if err := l.validator.Validate(in); err != nil {
		l.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	acHash, err := l.secHash.Hash(in.AccessToken)
	if err != nil {
		l.tel.Logger().Error(ctx, "failed to hash access token", err)

		return nil, goerror.NewServerInternal(err)
	}

	if err := l.store.DeleteTokenByAccess(ctx, string(acHash)); err != nil {
		l.tel.Logger().Error(ctx, "failed to delete token by access token", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.LogoutOutput{Message: "You have successfully logged out!"}, nil
}
