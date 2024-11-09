package telemetry

import (
	"context"
	"log"
	"time"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func WithZapLogger(level logger.Level, filters []string) func(*Telemetry) {
	return func(t *Telemetry) {
		lo, err := logger.NewZapLogger(
			logger.ZapWithLevel(level),
			logger.ZapWithVerbose(true),
			logger.ZapWithFilteredKeys(filters),
		)
		if err != nil {
			log.Printf("error while initialize zap logger %v", err)

			return
		}

		t.logger = lo
		t.flushers = append(t.flushers, lo.Close)
	}
}

func WithConsoleTracer(serviceName string) func(*Telemetry) {
	return func(t *Telemetry) {
		traceExporter, _ := stdouttrace.New() //nolint:errcheck // it will never error
		tp := trace.NewTracerProvider(
			trace.WithBatcher(
				traceExporter,
				trace.WithBatchTimeout(time.Second),
			),
			trace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceName(serviceName),
				),
			),
		)

		t.tracer = tp
		t.tracerCollector = OPENTELEMETRY
		t.flushers = append(t.flushers, func() error {
			return tp.Shutdown(context.Background())
		})
	}
}

func WithOTLPTracer(address, serviceName string) func(*Telemetry) {
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
					semconv.ServiceNameKey.String(serviceName),
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
