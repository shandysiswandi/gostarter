package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey struct{}

type Claim struct {
	now   time.Time
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (c *Claim) Now() time.Time {
	return c.now
}

func NewClaim(email string, exp time.Duration, aud []string) *Claim {
	now := time.Now()

	return &Claim{
		now:   now,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gostarter",
			Subject:   email,
			Audience:  aud,
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

func ExtractClaimFromToken(token string) *Claim {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}

	claim := Claim{}
	if err := json.Unmarshal(payload, &claim); err != nil {
		return nil
	}

	return &claim
}

func SetClaimToContext(ctx context.Context, clm *Claim) context.Context {
	return context.WithValue(ctx, contextKey{}, clm)
}

func GetClaimFromContext(ctx context.Context) *Claim {
	token, ok := ctx.Value(contextKey{}).(*Claim)
	if !ok {
		return nil
	}

	return token
}
