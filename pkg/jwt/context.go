package jwt

import (
	"context"
)

type contextKey struct{}

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
