package interceptor

import (
	"context"
	"errors"
	"log"

	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
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
}
