package service

import (
	"context"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type ResetPasswordStore interface {
	FindPasswordResetByToken(ctx context.Context, t string) (*domain.PasswordReset, error)
	DeletePasswordReset(ctx context.Context, id uint64) error
	UpdateUserPassword(ctx context.Context, id uint64, pass string) error
}

type ResetPassword struct {
	telemetry *telemetry.Telemetry
	store     ResetPasswordStore
	validator validation.Validator
	hash      hash.Hash
	now       func() time.Time
}

func NewResetPassword(t *telemetry.Telemetry, v validation.Validator, h hash.Hash,
	s ResetPasswordStore,
) *ResetPassword {
	return &ResetPassword{
		telemetry: t,
		validator: v,
		hash:      h,
		store:     s,
		now:       time.Now,
	}
}

func (s *ResetPassword) Call(ctx context.Context, in domain.ResetPasswordInput) (
	*domain.ResetPasswordOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "ResetPassword")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	ps, err := s.store.FindPasswordResetByToken(ctx, in.Token)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get password reset", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	if ps == nil {
		s.telemetry.Logger().Warn(ctx, "password reset not found")

		return nil, goerror.NewBusiness("invalid token", goerror.CodeUnauthorized)
	}

	if ps.ExpiresAt.Before(s.now()) {
		s.telemetry.Logger().Warn(ctx, "password reset token has expired")

		return nil, goerror.NewBusiness("token has expired", goerror.CodeUnauthorized)
	}

	if err := s.store.DeletePasswordReset(ctx, ps.ID); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to delete password reset", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	passHash, err := s.hash.Hash(in.Password)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to hash password", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	if err := s.store.UpdateUserPassword(ctx, ps.UserID, string(passHash)); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to delete password reset", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	return &domain.ResetPasswordOutput{}, nil
}
