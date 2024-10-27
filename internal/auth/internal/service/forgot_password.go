package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type ForgotPasswordStore interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindPasswordResetByUserID(ctx context.Context, uid uint64) (*domain.PasswordReset, error)
	SavePasswordReset(ctx context.Context, token domain.PasswordReset) error
	DeletePasswordReset(ctx context.Context, id uint64) error
}

type ForgotPassword struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	idnum     uid.NumberID
	secHash   hash.Hash
	store     ForgotPasswordStore
	now       func() time.Time
}

func NewForgotPassword(t *telemetry.Telemetry, v validation.Validator, idnum uid.NumberID,
	secHash hash.Hash, s ForgotPasswordStore,
) *ForgotPassword {
	return &ForgotPassword{
		telemetry: t,
		validator: v,
		idnum:     idnum,
		secHash:   secHash,
		store:     s,
		now:       time.Now,
	}
}

func (s *ForgotPassword) Call(ctx context.Context, in domain.ForgotPasswordInput) (
	*domain.ForgotPasswordOutput, error,
) {
	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	user, err := s.store.FindUserByEmail(ctx, in.Email)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get user", err, logger.String("email", in.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	if user == nil {
		s.telemetry.Logger().Warn(ctx, "user not found", logger.String("email", in.Email))

		return &domain.ForgotPasswordOutput{}, nil
	}

	ps, err := s.store.FindPasswordResetByUserID(ctx, user.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get password reset", err, logger.String("email", in.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	return s.doBest(ctx, in, user, ps)
}

func (s *ForgotPassword) doBest(ctx context.Context, in domain.ForgotPasswordInput, user *domain.User,
	ps *domain.PasswordReset,
) (*domain.ForgotPasswordOutput, error) {
	now := s.now()
	if ps != nil {
		if !ps.ExpiresAt.Before(now) {
			return &domain.ForgotPasswordOutput{}, nil
		}

		if err := s.store.DeletePasswordReset(ctx, ps.ID); err != nil {
			s.telemetry.Logger().Error(ctx, "failed to delete password reset", err,
				logger.String("email", in.Email))

			return nil, goerror.NewServer("internal server error", err)
		}
	}

	token, err := s.secHash.Hash(fmt.Sprintf("%d-%v", user.ID, now.Unix()))
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to generate password reset token", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	psData := domain.PasswordReset{
		ID:        s.idnum.Generate(),
		UserID:    user.ID,
		Token:     string(token),
		ExpiresAt: now.Add(time.Hour),
	}
	if err := s.store.SavePasswordReset(ctx, psData); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to save password reset", err,
			logger.String("email", in.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	return &domain.ForgotPasswordOutput{}, nil
}
