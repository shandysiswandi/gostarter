package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

const msgSuccess = "If an account with this email exists, you'll receive a password reset email shortly."

type ForgotPasswordStore interface {
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
	PasswordResetByUserID(ctx context.Context, uid uint64) (*domain.PasswordReset, error)
	PasswordResetSave(ctx context.Context, ps domain.PasswordReset) error
	PasswordResetDelete(ctx context.Context, id uint64) error
}

type ForgotPassword struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	idnum     uid.NumberID
	secHash   hash.Hash
	clock     clock.Clocker
	store     ForgotPasswordStore
}

func NewForgotPassword(dep Dependency, s ForgotPasswordStore) *ForgotPassword {
	return &ForgotPassword{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
		idnum:     dep.UIDNumber,
		secHash:   dep.SecHash,
		clock:     dep.Clock,
		store:     s,
	}
}

func (s *ForgotPassword) Call(ctx context.Context, in domain.ForgotPasswordInput) (
	*domain.ForgotPasswordOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.usecase.ForgotPassword")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	user, err := s.store.UserByEmail(ctx, in.Email)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		s.telemetry.Logger().Warn(ctx, "user not found", logger.KeyVal("email", in.Email))

		return &domain.ForgotPasswordOutput{
			Email:   in.Email,
			Message: msgSuccess,
		}, nil
	}

	ps, err := s.store.PasswordResetByUserID(ctx, user.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get password reset", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	return s.processPasswordReset(ctx, in, user, ps)
}

func (s *ForgotPassword) processPasswordReset(ctx context.Context, in domain.ForgotPasswordInput,
	user *domain.User, ps *domain.PasswordReset,
) (*domain.ForgotPasswordOutput, error) {
	now := s.clock.Now()
	if ps != nil {
		if !ps.ExpiresAt.Before(now) {
			return &domain.ForgotPasswordOutput{
				Email:   in.Email,
				Message: msgSuccess,
			}, nil
		}

		if err := s.store.PasswordResetDelete(ctx, ps.ID); err != nil {
			s.telemetry.Logger().Error(ctx, "failed to delete password reset", err,
				logger.KeyVal("email", in.Email))

			return nil, goerror.NewServerInternal(err)
		}
	}

	token, err := s.secHash.Hash(fmt.Sprintf("%d-%v", user.ID, now.Unix()))
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to generate password reset token", err)

		return nil, goerror.NewServerInternal(err)
	}

	psData := domain.PasswordReset{
		ID:        s.idnum.Generate(),
		UserID:    user.ID,
		Token:     string(token),
		ExpiresAt: now.Add(time.Hour),
	}
	if err := s.store.PasswordResetSave(ctx, psData); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to save password reset", err,
			logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	return &domain.ForgotPasswordOutput{
		Email:   in.Email,
		Message: msgSuccess,
	}, nil
}
