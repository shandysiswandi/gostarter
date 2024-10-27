package service

import (
	"context"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
)

type generateAndSaveTokenOutput struct {
	accessToken      string
	refreshToken     string
	accessExpiresIn  int64 // in seconds
	refreshExpiresIn int64 // in seconds
}

func generateAndSaveToken(ctx context.Context, log logger.Logger, jwte jwt.JWT, secHash hash.Hash,
	store func(context.Context, domain.Token) error, tid, uid uint64, email string) (
	*generateAndSaveTokenOutput, error,
) {
	acClaim := jwt.NewClaim(email, time.Hour, []string{"gostarter.access.token"})
	accToken, err := jwte.Generate(acClaim)
	if err != nil {
		log.Error(ctx, "failed to generate access token", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	refClaim := jwt.NewClaim(email, time.Hour*24, []string{"gostarter.refresh.token"})
	refToken, err := jwte.Generate(refClaim)
	if err != nil {
		log.Error(ctx, "failed to generate refresh token", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	acHash, err := secHash.Hash(accToken)
	if err != nil {
		log.Error(ctx, "failed to hash access token", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	refHash, err := secHash.Hash(refToken)
	if err != nil {
		log.Error(ctx, "failed to hash refresh token", err)

		return nil, goerror.NewServer("internal server error", err)
	}

	token := domain.Token{
		ID:               tid,
		UserID:           uid,
		AccessToken:      string(acHash),
		RefreshToken:     string(refHash),
		AccessExpiredAt:  acClaim.ExpiresAt.Time,
		RefreshExpiredAt: refClaim.ExpiresAt.Time,
	}
	if err := store(ctx, token); err != nil {
		log.Error(ctx, "failed to save tokens", err, logger.String("email", email))

		return nil, goerror.NewServer("internal server error", err)
	}

	return &generateAndSaveTokenOutput{
		accessToken:      accToken,
		refreshToken:     refToken,
		accessExpiresIn:  acClaim.ExpiresAt.Time.Unix() - acClaim.Now().Unix(),
		refreshExpiresIn: refClaim.ExpiresAt.Time.Unix() - refClaim.Now().Unix(),
	}, nil
}
