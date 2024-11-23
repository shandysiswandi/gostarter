package telemetry

import (
	"errors"
	"log"
	"os"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/filter"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Collector int

const (
	NOOP Collector = iota
	OPENTELEMETRY
)

type Telemetry struct {
	name            string
	filter          *filter.Filter
	logger          logger.Logger
	tracer          trace.TracerProvider
	tracerCollector Collector
	flushers        []func() error
}

func NewTelemetry(opts ...func(*Telemetry)) *Telemetry {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	tel := &Telemetry{
		name:            "telemetry",
		filter:          filter.NewFilter(),
		logger:          logger.NewNoopLogger(),
		tracer:          noop.NewTracerProvider(),
		tracerCollector: NOOP,
		flushers:        make([]func() error, 0),
	}

	for _, opt := range opts {
		opt(tel)
	}

	return tel
}

func (t *Telemetry) Close() error {
	var errs error

	for _, flusher := range t.flushers {
		if err := flusher(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (t *Telemetry) Logger() logger.Logger {
	return t.logger
}

func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer.Tracer(t.name)
}

func (t *Telemetry) TracerProvider() trace.TracerProvider {
	return t.tracer
}

func (t *Telemetry) TracerCollector() Collector {
	return t.tracerCollector
}

func (t *Telemetry) Meter() metric.Meter {
	return nil
}
