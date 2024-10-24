package telemetry

import (
	"errors"
	"log"
	"os"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Collector int

const (
	NOOP Collector = iota
	OPENTELEMETRY
)

type Telemetry struct {
	name     string
	logger   logger.Logger
	flushers []func() error
}

func NewTelemetry(opts ...func(*Telemetry)) *Telemetry {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	tel := &Telemetry{
		name:     "telemetry",
		logger:   logger.NewNoopLogger(),
		flushers: make([]func() error, 0),
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
	return nil
}

func (t *Telemetry) Meter() metric.Meter {
	return nil
}

func WithZapLogger(level logger.Level) func(*Telemetry) {
	return func(t *Telemetry) {
		l, err := logger.NewZapLogger(level)
		if err != nil && err != os.ErrInvalid {
			log.Printf("error while initialize zap logger %v", err)
			return
		}

		t.logger = l
		t.flushers = append(t.flushers, l.Close)
	}
}
