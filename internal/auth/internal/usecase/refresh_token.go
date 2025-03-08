package usecase

import (
	"context"
	"time"

	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/jwt"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/lib"
)

type RefreshTokenStore interface {
	TokenByRefresh(ctx context.Context, ref string) (*domain.Token, error)
	TokenSave(ctx context.Context, token domain.Token) error
	TokenUpdate(ctx context.Context, token domain.Token) error
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
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.usecase.RefreshToken")
	defer span.End()

	if err := s.validator.Validate(in); err != nil {
		s.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	refHash, err := s.secHash.Hash(in.RefreshToken)
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to hash refresh token", err)

		return nil, goerror.NewServerInternal(err)
	}

	refToken, err := s.store.TokenByRefresh(ctx, string(refHash))
	if err != nil {
		s.telemetry.Logger().Error(ctx, "failed to get token", err,
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewServerInternal(err)
	}

	if refToken == nil {
		s.telemetry.Logger().Warn(ctx, "token not found",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	if refToken.RefreshExpiredAt.Before(time.Now()) {
		s.telemetry.Logger().Warn(ctx, "token has expired",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("Token has expired", goerror.CodeUnauthorized)
	}

	clm := lib.ExtractJWTClaim(in.RefreshToken)
	if clm == nil {
		s.telemetry.Logger().Warn(ctx, "token is malformed",
			logger.KeyVal("refresh_token_hash", string(refHash)))

		return nil, goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized)
	}

	tgsIn := tokenGenSaverIn{email: clm.Subject, token: refToken, userID: refToken.UserID}
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
