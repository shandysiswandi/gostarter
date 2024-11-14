package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type RegisterStore interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveUser(ctx context.Context, token domain.User) error
}

type Register struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	hash      hash.Hash
	store     RegisterStore
}

func NewRegister(t *telemetry.Telemetry, v validation.Validator,
	idnum uid.NumberID, hash hash.Hash, s RegisterStore,
) *Register {
	return &Register{
		telemetry: t,
		validator: v,
		uidnumber: idnum,
		hash:      hash,
		store:     s,
	}
}

func (s *Register) Call(ctx context.Context, in domain.RegisterInput) (*domain.RegisterOutput, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "service.Register")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	user, err := s.store.FindUserByEmail(ctx, in.Email)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	if user != nil {
		s.telemetry.Logger().Warn(ctx, "user already exists", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("email already registered", goerror.CodeConflict)
	}

	passHash, err := s.hash.Hash(in.Password)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to hash password", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	userData := domain.User{
		ID:       s.uidnumber.Generate(),
		Email:    in.Email,
		Password: string(passHash),
	}
	if err := s.store.SaveUser(ctx, userData); err != nil {
		s.telemetry.Logger().Error(ctx, "failed to save user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	return &domain.RegisterOutput{}, nil
}
