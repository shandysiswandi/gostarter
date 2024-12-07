package telemetry

import (
	"errors"
	"log"
	"os"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/filter"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"go.opentelemetry.io/otel/metric"
	mnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	tnoop "go.opentelemetry.io/otel/trace/noop"
)

type Vendor int

const (
	NOOP Vendor = iota
	OPENTELEMETRY
)

type Telemetry struct {
	name      string
	verbose   bool
	filter    *filter.Filter
	logger    logger.Logger
	tracer    trace.TracerProvider
	meter     metric.MeterProvider
	collector Vendor
	flushers  []func() error
}

func NewTelemetry(opts ...Option) *Telemetry {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	tel := &Telemetry{
		name:      "telemetry",
		filter:    nil,
		verbose:   false,
		logger:    logger.NewNoopLogger(),
		tracer:    tnoop.NewTracerProvider(),
		meter:     mnoop.NewMeterProvider(),
		collector: NOOP,
		flushers:  make([]func() error, 0),
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

func (t *Telemetry) Collector() Vendor {
	return t.collector
}

func (t *Telemetry) Meter() metric.Meter {
	return t.meter.Meter(t.name)
}

func (t *Telemetry) Filter() *filter.Filter {
	return t.filter
}

func (t *Telemetry) Verbose() bool {
	return t.verbose
}
