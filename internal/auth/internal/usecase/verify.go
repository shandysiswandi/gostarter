package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type VerifyStore interface {
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
	UserVerificationByUserID(ctx context.Context, uid uint64) (*domain.UserVerification, error)
}

type Verify struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	secHash   hash.Hash
	store     VerifyStore
}

func NewVerify(dep Dependency, s VerifyStore) *Verify {
	return &Verify{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		secHash:   dep.Hash,
		store:     s,
	}
}

func (s *Verify) Call(ctx context.Context, in domain.VerifyInput) (*domain.VerifyOutput, error) {
	ctx, span := s.tel.Tracer().Start(ctx, "auth.usecase.Verify")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	user, err := s.store.UserByEmail(ctx, in.Email)
	if err != nil {
		s.tel.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	if user == nil {
		s.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	if user.VerifiedAt.Valid {
		return &domain.VerifyOutput{
			Email:    user.Email,
			VerifyAt: user.VerifiedAt.V,
		}, nil
	}

	// find user and code
	uv, err := s.store.UserVerificationByUserID(ctx, user.ID)
	if err != nil {
		s.tel.Logger().Error(ctx, "failed to get user verification", err, logger.KeyVal("user.id", user.ID))

		return nil, goerror.NewServerInternal(err)
	}

	if uv == nil {
		s.tel.Logger().Warn(ctx, "user verification not found", logger.KeyVal("user.id", user.ID))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	// update user verify

	return nil, nil
}
