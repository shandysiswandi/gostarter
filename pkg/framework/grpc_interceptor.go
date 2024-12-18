package framework

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerRecovery(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (
	_ any, err error,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic because: %v\n", r)
			debug.PrintStack()

			err = status.Error(codes.Internal, "internal server error")
		}
	}()

	return next(ctx, req)
}

func UnaryServerError(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (
	any, error,
) {
	resp, err := next(ctx, req)
	if err != nil {
		var errs *goerror.GoError
		if ok := errors.As(err, &errs); ok {
			dataMD := make(metadata.MD)
			dataMD.Set("error_code", errs.Code().String())
			dataMD.Set("error_type", errs.Type().String())
			dataMD.Set("error_message", errs.Msg())

			if err := grpc.SetTrailer(ctx, dataMD); err != nil {
				log.Print(err)
			}

			return nil, errs.GRPCStatus().Err()
		}

		return nil, err
	}

	return resp, nil
}

func UnaryServerJWT(audience string, skipMethods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if len(skipMethods) > 0 {
			for _, prefix := range skipMethods {
				if strings.HasPrefix(info.FullMethod, prefix) {
					return next(ctx, req)
				}
			}
		}

		return doUnaryServerJWT(ctx, req, next, audience)
	}
}

func doUnaryServerJWT(ctx context.Context, req any, next grpc.UnaryHandler,
	audience string,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, goerror.NewServerInternal(nil)
	}

	clm := jwt.ExtractClaimFromToken(strings.TrimPrefix(md["authorization"][0], "Bearer "))
	if !clm.VerifyAudience(audience, true) {
		return nil, goerror.NewBusiness("invalid token audience", goerror.CodeUnauthorized)
	}

	return next(jwt.SetClaim(ctx, clm), req)
}

func UnaryServerProtoValidate(validator validation.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := validator.Validate(req); err != nil {
			return nil, goerror.NewInvalidInput("proto validation failed", err)
		}

		return next(ctx, req)
	}
}
