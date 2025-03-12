package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
)

type UpdatePasswordStore interface {
	User(ctx context.Context, id uint64) (*domain.User, error)
	UserUpdate(ctx context.Context, user domain.User) error
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

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	var uid uint64
	if clm := lib.GetJWTClaim(ctx); clm != nil {
		uid = clm.AuthID
	}

	user, err := up.store.User(ctx, uid)
	if err != nil {
		up.tel.Logger().Error(ctx, "failed to find user", err, logger.KeyVal("id", uid))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		up.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("id", uid))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	if !up.hash.Verify(user.Password, in.CurrentPassword) {
		up.tel.Logger().Warn(ctx, "password not match", logger.KeyVal("id", uid))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
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
