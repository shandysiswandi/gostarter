package usecase

import (
	"context"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type ResetPasswordStore interface {
	PasswordResetByToken(ctx context.Context, t string) (*domain.PasswordReset, error)
	PasswordResetDelete(ctx context.Context, id uint64) error
	UserUpdatePassword(ctx context.Context, id uint64, pass string) error
}

type ResetPassword struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	hash      hash.Hash
	store     ResetPasswordStore
}

func NewResetPassword(dep Dependency, s ResetPasswordStore) *ResetPassword {
	return &ResetPassword{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		hash:      dep.Hash,
		store:     s,
	}
}

func (s *ResetPassword) Call(ctx context.Context, in domain.ResetPasswordInput) (
	*domain.ResetPasswordOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.usecase.ResetPassword")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	ps, err := s.store.PasswordResetByToken(ctx, in.Token)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get password reset", err)

		return nil, goerror.NewServerInternal(err)
	}

	if ps == nil {
		s.telemetry.Logger().Warn(ctx, "password reset not found")

		return nil, goerror.NewBusiness("Invalid token", goerror.CodeUnauthorized)
	}

	if ps.ExpiresAt.Before(time.Now()) {
		s.telemetry.Logger().Warn(ctx, "password reset token has expired")

		return nil, goerror.NewBusiness("Token has expired", goerror.CodeUnauthorized)
	}

	if err := s.store.PasswordResetDelete(ctx, ps.ID); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to delete password reset", err)

		return nil, goerror.NewServerInternal(err)
	}

	passHash, err := s.hash.Hash(in.Password)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to hash password", err)

		return nil, goerror.NewServerInternal(err)
	}

	if err := s.store.UserUpdatePassword(ctx, ps.UserID, string(passHash)); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to delete password reset", err)

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.ResetPasswordOutput{
		Message: "Your password has been successfully reset.",
	}, nil
}
