package usecase

import (
	"context"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type RefreshTokenStore interface {
	FindTokenByRefresh(ctx context.Context, ref string) (*domain.Token, error)
	SaveToken(ctx context.Context, token domain.Token) error
	UpdateToken(ctx context.Context, token domain.Token) error
}

type RefreshToken struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	secHash   hash.Hash
	jwt       jwt.JWT
	clock     clock.Clocker
	store     RefreshTokenStore
	tgs       *tokenGenSaver
}

func NewRefreshToken(dep Dependency, s RefreshTokenStore) *RefreshToken {
	return &RefreshToken{
		telemetry: dep.Telemetry,
		validator: dep.Validator,
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

func (s *RefreshToken) Call(ctx context.Context, in domain.RefreshTokenInput) (
	*domain.RefreshTokenOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "usecase.RefreshToken")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	refHash, err := s.secHash.Hash(in.RefreshToken)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to hash refresh token", err)

		return nil, goerror.NewServerInternal(err)
	}

	refToken, err := s.store.FindTokenByRefresh(ctx, string(refHash))
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get token", err,
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewServerInternal(err)
	}

	if refToken == nil {
		s.telemetry.Logger().Warn(ctx, "token not found",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	if refToken.RefreshExpiredAt.Before(time.Now()) {
		s.telemetry.Logger().Warn(ctx, "token has expired",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("token has expired", goerror.CodeUnauthorized)
	}

	clm := jwt.ExtractClaimFromToken(in.RefreshToken)
	if clm == nil {
		s.telemetry.Logger().Warn(ctx, "token is malformed",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("invalid credentials", goerror.CodeUnauthorized)
	}

	tgsIn := tokenGenSaverIn{email: clm.Email, token: refToken, userID: refToken.UserID}
	tgso, err := s.tgs.do(ctx, tgsIn)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenOutput{
		AccessToken:      tgso.accessToken,
		RefreshToken:     tgso.refreshToken,
		AccessExpiresIn:  tgso.accessExpiresIn,
		RefreshExpiresIn: tgso.refreshExpiresIn,
	}, nil
}
