package interceptor

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (
		_ any, err error,
	) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic because: %v\n", r)
				debug.PrintStack()

				err = status.Error(codes.Internal, "Internal Server Error")
			}
		}()

		return next(ctx, req)
	}
}
