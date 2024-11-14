package interceptor

import (
	"context"

	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
)

func UnaryServerProtoValidate(validator validation.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := validator.Validate(req); err != nil {
			return nil, goerror.NewInvalidInput("proto validation failed", err)
		}

		return next(ctx, req)
	}
}
