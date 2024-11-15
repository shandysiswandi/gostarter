package service

import (
	"context"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type RefreshTokenStore interface {
	FindTokenByRefresh(ctx context.Context, ref string) (*domain.Token, error)
	SaveToken(ctx context.Context, token domain.Token) error
}

type RefreshToken struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	secHash   hash.Hash
	jwt       jwt.JWT
	store     RefreshTokenStore
}

func NewRefreshToken(t *telemetry.Telemetry, v validation.Validator,
	idnum uid.NumberID, secHash hash.Hash, j jwt.JWT, s RefreshTokenStore,
) *RefreshToken {
	return &RefreshToken{
		telemetry: t,
		validator: v,
		uidnumber: idnum,
		secHash:   secHash,
		jwt:       j,
		store:     s,
	}
}

func (s *RefreshToken) Call(ctx context.Context, in domain.RefreshTokenInput) (
	*domain.RefreshTokenOutput, error,
) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "service.RefreshToken")
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

	resp, err := generateAndSaveToken(ctx, s.telemetry.Logger(), s.jwt, s.secHash, s.store.SaveToken,
		refToken.ID, refToken.UserID, clm.Email)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenOutput{
		AccessToken:      resp.accessToken,
		RefreshToken:     resp.refreshToken,
		AccessExpiresIn:  resp.accessExpiresIn,
		RefreshExpiresIn: resp.refreshExpiresIn,
	}, nil
}
