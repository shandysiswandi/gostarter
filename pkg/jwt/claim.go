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
	AuthID uint64 `json:"auth_id,string"`
	jwt.RegisteredClaims
}

func NewClaim(authID uint64, email string, exp time.Time, aud []string) *Claim {
	return &Claim{
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

func SetClaim(ctx context.Context, clm *Claim) context.Context {
	return context.WithValue(ctx, contextKey{}, clm)
}

func GetClaim(ctx context.Context) *Claim {
	token, ok := ctx.Value(contextKey{}).(*Claim)
	if !ok {
		return nil
	}

	return token
}
