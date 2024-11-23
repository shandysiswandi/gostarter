package telemetry

import (
	"context"
	"log"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/filter"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func WithServiceName(serviceName string) func(*Telemetry) {
	return func(t *Telemetry) {
		t.name = serviceName
	}
}

func WithLogFilter(keys ...string) func(*Telemetry) {
	return func(t *Telemetry) {
		t.filter = filter.NewFilter(
			filter.WithHeaders(keys...),
		)
	}
}

func WithZapLogger(lvl logger.Level, fKeys ...string) func(*Telemetry) {
	return func(t *Telemetry) {
		lo, err := logger.NewZapLogger(lvl, fKeys...)
		if err != nil {
			log.Printf("error while initialize zap logger %v", err)

			return
		}

		t.logger = lo
		t.flushers = append(t.flushers, lo.Close)
	}
}

func WithOTLPTracer(address string) func(*Telemetry) {
	return func(t *Telemetry) {
		ctx := context.Background()

		exporter, err := otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(address),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			log.Printf("error while initialize open telemetry tracer %v", err)

			return
		}

		tp := trace.NewTracerProvider(
			trace.WithBatcher(exporter),
			trace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(t.name),
				),
			),
		)

		t.tracer = tp
		t.tracerCollector = OPENTELEMETRY
		t.flushers = append(t.flushers, func() error {
			return tp.Shutdown(ctx)
		})
	}
}
