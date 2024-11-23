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
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Dependency struct {
	Telemetry *telemetry.Telemetry
	Validator validation.Validator
	UIDNumber uid.NumberID
	Hash      hash.Hash
	SecHash   hash.Hash
	JWT       jwt.JWT
	Clock     clock.Clocker
}

type tokenSaver interface {
	SaveToken(ctx context.Context, token domain.Token) error
}

type tokenGenSaver struct {
	jwt     jwt.JWT
	tel     *telemetry.Telemetry
	secHash hash.Hash
	clock   clock.Clocker
	ts      tokenSaver
}

type tokenGenSaverIn struct {
	email   string
	tokenID uint64
	userID  uint64
}

type tokenGenSaverOut struct {
	accessToken      string
	refreshToken     string
	accessExpiresIn  int64 // in seconds
	refreshExpiresIn int64 // in seconds
}

func (tgs *tokenGenSaver) do(ctx context.Context, in tokenGenSaverIn) (*tokenGenSaverOut, error) {
	now := tgs.clock.Now()

	acClaim := jwt.NewClaim(in.email, time.Hour, now, []string{"gostarter.access.token"})
	accToken, err := tgs.jwt.Generate(acClaim)
	if err != nil {
		tgs.tel.Logger().Error(ctx, "failed to generate access token", err)

		return nil, goerror.NewServerInternal(err)
	}

	refClaim := jwt.NewClaim(in.email, time.Hour*24, now, []string{"gostarter.refresh.token"})
	refToken, err := tgs.jwt.Generate(refClaim)
	if err != nil {
		tgs.tel.Logger().Error(ctx, "failed to generate refresh token", err)

		return nil, goerror.NewServerInternal(err)
	}

	acHash, err := tgs.secHash.Hash(accToken)
	if err != nil {
		tgs.tel.Logger().Error(ctx, "failed to hash access token", err)

		return nil, goerror.NewServerInternal(err)
	}

	refHash, err := tgs.secHash.Hash(refToken)
	if err != nil {
		tgs.tel.Logger().Error(ctx, "failed to hash refresh token", err)

		return nil, goerror.NewServerInternal(err)
	}

	token := domain.Token{
		ID:               in.tokenID,
		UserID:           in.userID,
		AccessToken:      string(acHash),
		RefreshToken:     string(refHash),
		AccessExpiredAt:  acClaim.ExpiresAt.Time,
		RefreshExpiredAt: refClaim.ExpiresAt.Time,
	}
	if err := tgs.ts.SaveToken(ctx, token); err != nil {
		tgs.tel.Logger().Error(ctx, "failed to save tokens", err, logger.KeyVal("email", in.email))

		return nil, goerror.NewServerInternal(err)
	}

	return &tokenGenSaverOut{
		accessToken:      accToken,
		refreshToken:     refToken,
		accessExpiresIn:  acClaim.ExpiresAt.Time.Unix() - now.Unix(),
		refreshExpiresIn: refClaim.ExpiresAt.Time.Unix() - now.Unix(),
	}, nil
}
