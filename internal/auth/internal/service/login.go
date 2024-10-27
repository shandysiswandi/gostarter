package service

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type LoginStore interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindTokenByUserID(ctx context.Context, uid uint64) (*domain.Token, error)
	SaveToken(ctx context.Context, token domain.Token) error
}

type Login struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	hash      hash.Hash
	secHash   hash.Hash
	jwt       jwt.JWT
	store     LoginStore
}

func NewLogin(t *telemetry.Telemetry, v validation.Validator, idnum uid.NumberID, hash, secHash hash.Hash,
	j jwt.JWT, s LoginStore) *Login {
	return &Login{
		telemetry: t,
		validator: v,
		uidnumber: idnum,
		hash:      hash,
		secHash:   secHash,
		jwt:       j,
		store:     s,
	}
}

func (s *Login) Call(ctx context.Context, in domain.LoginInput) (*domain.LoginOutput, error) {
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

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	if !s.hash.Verify(user.Password, in.Password) {
		s.telemetry.Logger().Warn(ctx, "password not match", logger.String("email", in.Email))

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	token, err := s.store.FindTokenByUserID(ctx, user.ID)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get token", err, logger.String("email", user.Email))

		return nil, goerror.NewServer("internal server error", err)
	}

	tid := s.uidnumber.Generate()
	if token != nil {
		tid = token.ID
	}

	resp, err := generateAndSaveToken(ctx, s.telemetry.Logger(), s.jwt, s.secHash, s.store.SaveToken,
		tid, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.LoginOutput{
		AccessToken:      resp.accessToken,
		RefreshToken:     resp.refreshToken,
		AccessExpiresIn:  resp.accessExpiresIn,
		RefreshExpiresIn: resp.refreshExpiresIn,
	}, nil
}
