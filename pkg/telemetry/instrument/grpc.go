package instrument

import (
	"context"

	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/requestid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryTelemetryServerInterceptor(tel *telemetry.Telemetry, sid func() string) []grpc.ServerOption {
	var opts []grpc.ServerOption

	ins := &grpcServer{uuid: sid, tel: tel}

	switch tel.TracerCollector() {
	case telemetry.OPENTELEMETRY:
		tracerProvider := tel.TracerProvider()

		statsHandler := otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tracerProvider),
		)

		opts = append(opts, grpc.UnaryInterceptor(ins.log), grpc.StatsHandler(statsHandler))

		return opts

	case telemetry.NOOP:
		return nil

	default:
		return nil
	}
}

type grpcServer struct {
	uuid func() string
	tel  *telemetry.Telemetry
}

func (ins *grpcServer) log(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (
	any, error,
) {
	ctx = ins.requestID(ctx)

	resp, err := h(ctx, req)

	ins.tel.Logger().Info(ctx, "grpc request response",
		logger.KeyVal("rpc.method", info.FullMethod),
		logger.KeyVal("rpc.status.code", int32(status.Code(err))),
		logger.KeyVal("rpc.status.name", status.Code(err).String()),
	)

	return resp, err
}

func (ins *grpcServer) requestID(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	if reqID, ok := md["x-request-id"]; ok && len(reqID) > 0 {
		ctx = requestid.Set(ctx, reqID[0])
	}

	return requestid.Set(ctx, ins.uuid())
}
