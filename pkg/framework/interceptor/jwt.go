package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerJWT(jwte jwt.JWT, audience string, skipMethods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if len(skipMethods) > 0 {
			for _, prefix := range skipMethods {
				if strings.HasPrefix(info.FullMethod, prefix) {
					return next(ctx, req)
				}
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, goerror.NewServer("internal server error", nil)
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return nil, goerror.NewBusiness("authorization header missing", goerror.CodeUnauthorized)
		}

		if !strings.HasPrefix(authHeader[0], "Bearer ") {
			return nil, goerror.NewBusiness("invalid format", goerror.CodeUnauthorized)
		}

		clm, err := jwte.Verify(strings.TrimPrefix(authHeader[0], "Bearer "))
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, goerror.NewBusiness("expired token", goerror.CodeUnauthorized)
		}

		if err != nil {
			return nil, goerror.NewBusiness("invalid token", goerror.CodeUnauthorized)
		}

		if !clm.VerifyAudience(audience, true) {
			return nil, goerror.NewBusiness("invalid token audience", goerror.CodeUnauthorized)
		}

		return next(jwt.SetClaim(ctx, clm), req)
	}
}
