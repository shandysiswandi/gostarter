package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type RegisterStore interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveUser(ctx context.Context, user domain.User) error
	SaveAccount(ctx context.Context, user domain.Account) error
}

type Register struct {
	tele      *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	hash      hash.Hash
	trx       dbops.Tx
	store     RegisterStore
}

func NewRegister(dep Dependency, s RegisterStore) *Register {
	return &Register{
		tele:      dep.Telemetry,
		validator: dep.Validator,
		uidnumber: dep.UIDNumber,
		hash:      dep.Hash,
		trx:       dep.Transaction,
		store:     s,
	}
}

func (s *Register) Call(ctx context.Context, in domain.RegisterInput) (*domain.RegisterOutput, error) {
	ctx, span := s.tele.Tracer().Start(ctx, "auth.usecase.Register")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.tele.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	user, err := s.store.FindUserByEmail(ctx, in.Email)
	if err != nil {
		s.tele.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	if user != nil {
		s.tele.Logger().Warn(ctx, "user already exists", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("Email already registered", goerror.CodeConflict)
	}

	passHash, err := s.hash.Hash(in.Password)
	if err != nil {
		s.tele.Logger().Error(ctx, "failed to hash password", err)

		return nil, goerror.NewServerInternal(err)
	}

	err = s.trx.Transaction(ctx, func(ctx context.Context) error {
		userData := domain.User{
			ID:       s.uidnumber.Generate(),
			Name:     in.Name,
			Email:    in.Email,
			Password: string(passHash),
		}
		if err := s.store.SaveUser(ctx, userData); err != nil {
			s.tele.Logger().Error(ctx, "failed to save user", err, logger.KeyVal("email", in.Email))

			return goerror.NewServerInternal(err)
		}

		accountData := domain.Account{ID: s.uidnumber.Generate(), UserID: userData.ID}
		if err := s.store.SaveAccount(ctx, accountData); err != nil {
			s.tele.Logger().Error(ctx, "failed to save account", err, logger.KeyVal("email", in.Email))

			return goerror.NewServerInternal(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &domain.RegisterOutput{Email: in.Email}, nil
}
