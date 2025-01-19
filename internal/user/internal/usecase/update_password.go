package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type UpdatePasswordStore interface {
	FindUser(ctx context.Context, id uint64) (*domain.User, error)
	UpdatePassword(ctx context.Context, user domain.User) error
}

type UpdatePassword struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	hash      hash.Hash
	store     UpdatePasswordStore
}

func NewUpdatePassword(dep Dependency, s UpdatePasswordStore) *UpdatePassword {
	return &UpdatePassword{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		hash:      dep.Hash,
		store:     s,
	}
}

func (up *UpdatePassword) Call(ctx context.Context, in domain.UpdatePasswordInput) (*domain.User, error) {
	ctx, span := up.tel.Tracer().Start(ctx, "user.usecase.UpdatePassword")
	defer span.End()

	if err := up.validator.Validate(in); err != nil {
		up.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	var uid uint64
	if clm := jwt.GetClaim(ctx); clm != nil {
		uid = clm.AuthID
	}

	user, err := up.store.FindUser(ctx, uid)
	if err != nil {
		up.tel.Logger().Error(ctx, "failed to find user", err, logger.KeyVal("id", uid))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		up.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("id", uid))

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	if !up.hash.Verify(user.Password, in.CurrentPassword) {
		up.tel.Logger().Warn(ctx, "password not match", logger.KeyVal("id", uid))

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	passHash, err := up.hash.Hash(in.NewPassword)
	if err != nil {
		up.tel.Logger().Error(ctx, "failed to hash password", err)

		return nil, goerror.NewServerInternal(err)
	}

	user.Password = string(passHash)
	if err := up.store.UpdatePassword(ctx, *user); err != nil {
		up.tel.Logger().Error(ctx, "failed to update password user", err, logger.KeyVal("id", uid))

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.User{
		ID:       uid,
		Name:     user.Name,
		Email:    user.Email,
		Password: "***",
	}, nil
}
