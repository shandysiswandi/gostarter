package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/jwt"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type LoginStore interface {
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
	TokenByUserID(ctx context.Context, uid uint64) (*domain.Token, error)
	TokenSave(ctx context.Context, token domain.Token) error
	TokenUpdate(ctx context.Context, token domain.Token) error
}

type Login struct {
	tel       *telemetry.Telemetry
	validator validation.Validator
	hash      hash.Hash
	secHash   hash.Hash
	jwt       jwt.JWT
	clock     clock.Clocker
	store     LoginStore
	tgs       *tokenGenSaver
}

func NewLogin(dep Dependency, s LoginStore) *Login {
	return &Login{
		tel:       dep.Telemetry,
		validator: dep.Validator,
		hash:      dep.Hash,
		secHash:   dep.SecHash,
		jwt:       dep.JWT,
		clock:     dep.Clock,
		store:     s,
		tgs: &tokenGenSaver{
			uidnumber: dep.UIDNumber,
			jwt:       dep.JWT,
			tel:       dep.Telemetry,
			secHash:   dep.SecHash,
			clock:     dep.Clock,
			ts:        s,
		},
	}
}

func (s *Login) Call(ctx context.Context, in domain.LoginInput) (*domain.LoginOutput, error) {
	ctx, span := s.tel.Tracer().Start(ctx, "auth.usecase.Login")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.tel.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	u, err := s.store.UserByEmail(ctx, in.Email)
	if err != nil {
		s.tel.Logger().Error(ctx, "failed to get user", err, logger.KeyVal("email", in.Email))

		return nil, goerror.NewServerInternal(err)
	}

	if u == nil {
		s.tel.Logger().Warn(ctx, "user not found", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	if !u.VerifiedAt.Valid {
		s.tel.Logger().Warn(ctx, "user not verified", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("Account isn't verified, please check email", goerror.CodeForbidden)
	}

	if !s.hash.Verify(u.Password, in.Password) {
		s.tel.Logger().Warn(ctx, "password not match", logger.KeyVal("email", in.Email))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	token, err := s.store.TokenByUserID(ctx, u.ID)
	if err != nil {
		s.tel.Logger().Error(ctx, "failed to get token", err, logger.KeyVal("email", u.Email))

		return nil, goerror.NewServerInternal(err)
	}

	tgsIn := tokenGenSaverIn{email: in.Email, userID: u.ID, token: token}
	tgso, err := s.tgs.do(ctx, tgsIn)
	if err != nil {
		return nil, err
	}

	return &domain.LoginOutput{
		AccessToken:      tgso.accessToken,
		RefreshToken:     tgso.refreshToken,
		AccessExpiresIn:  tgso.accessExpiresIn,
		RefreshExpiresIn: tgso.refreshExpiresIn,
	}, nil
}
