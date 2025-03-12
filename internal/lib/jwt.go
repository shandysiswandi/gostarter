package lib

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextJWTKey struct{}

type JWTClaim struct {
	AuthID uint64 `json:"auth_id,string"`
	jwt.RegisteredClaims
}

func (c JWTClaim) Validate() error {

	return nil
}

func NewJWTClaim(authID uint64, email string, exp time.Time, aud []string) *JWTClaim {
	return &JWTClaim{
		AuthID: authID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "GO_STARTER",
			Subject:   email,
			Audience:  aud,
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func ExtractJWTClaim(token string) *JWTClaim {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}

	claim := JWTClaim{}
	if err := json.Unmarshal(payload, &claim); err != nil {
		return nil
	}

	return &claim
}

func SetJWTClaim(ctx context.Context, clm *JWTClaim) context.Context {
	return context.WithValue(ctx, contextJWTKey{}, clm)
}

func GetJWTClaim(ctx context.Context) *JWTClaim {
	token, ok := ctx.Value(contextJWTKey{}).(*JWTClaim)
	if !ok {
		return nil
	}

	return token
}
