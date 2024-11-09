package telemetry

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
)

func TestWithZapLogger(t *testing.T) {
	type args struct {
		level   logger.Level
		filters []string
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{
				level:   logger.InfoLevel,
				filters: []string{"token"},
			},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithZapLogger(a.level, a.filters)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tel := tt.mockFn(tt.args)

			assert.Len(t, tel.flushers, 1)
			assert.NotNil(t, tel.Logger())
		})
	}
}

func TestWithConsoleTracer(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{serviceName: "gostarter"},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithConsoleTracer(a.serviceName)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tel := tt.mockFn(tt.args)

			assert.Len(t, tel.flushers, 1)
			assert.NotNil(t, tel.Tracer())
		})
	}
}

func TestWithOTLPTracer(t *testing.T) {
	type args struct {
		address     string
		serviceName string
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{serviceName: "gostarter"},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithOTLPTracer(a.address, a.serviceName)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tel := tt.mockFn(tt.args)

			assert.Len(t, tel.flushers, 1)
			assert.NotNil(t, tel.Tracer())
		})
	}
}
