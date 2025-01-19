package telemetry

import (
	"errors"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/filter"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric"
	mnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	tnoop "go.opentelemetry.io/otel/trace/noop"
)

func TestNewTelemetry(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *Telemetry
	}{
		{
			name: "Success",
			args: args{
				opts: []Option{func(t *Telemetry) {}},
			},
			want: NewTelemetry(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewTelemetry(tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTelemetry_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *Telemetry
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
		{
			name:    "ErrorWithZapLogger",
			wantErr: errors.New("sync /dev/stdout: invalid argument"),
			mockFn: func() *Telemetry {
				return NewTelemetry(WithZapLogger("", logger.InfoLevel, false))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.mockFn().Close()
			if err != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestTelemetry_Logger(t *testing.T) {
	tests := []struct {
		name   string
		want   func() logger.Logger
		mockFn func() *Telemetry
	}{
		{
			name: "SuccessNoopLogger",
			want: func() logger.Logger {
				return &logger.NoopLogger{}
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := tt.mockFn().Logger()
			assert.Equal(t, tt.want(), logger)
		})
	}
}

func TestTelemetry_Tracer(t *testing.T) {
	tests := []struct {
		name   string
		want   func() trace.Tracer
		mockFn func() *Telemetry
	}{
		{
			name: "SuccessNoopTracer",
			want: func() trace.Tracer {
				return tnoop.NewTracerProvider().Tracer("telemetry")
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().Tracer()
			assert.Equal(t, tt.want(), tr)
		})
	}
}

func TestTelemetry_TracerProvider(t *testing.T) {
	tests := []struct {
		name   string
		want   func() trace.TracerProvider
		mockFn func() *Telemetry
	}{
		{
			name: "Success",
			want: func() trace.TracerProvider {
				return tnoop.NewTracerProvider()
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().TracerProvider()
			assert.Equal(t, tt.want(), tr)
		})
	}
}

func TestTelemetry_TracerCollector(t *testing.T) {
	tests := []struct {
		name   string
		want   func() Vendor
		mockFn func() *Telemetry
	}{
		{
			name: "Success",
			want: func() Vendor {
				return NOOP
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().Collector()
			assert.Equal(t, tt.want(), tr)
		})
	}
}

func TestTelemetry_Meter(t *testing.T) {
	tests := []struct {
		name   string
		want   func() metric.Meter
		mockFn func() *Telemetry
	}{
		{
			name: "SuccessNoopMeter",
			want: func() metric.Meter {
				return mnoop.NewMeterProvider().Meter("telemetry")
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().Meter()
			assert.Equal(t, tt.want(), tr)
		})
	}
}

func TestTelemetry_Filter(t *testing.T) {
	tests := []struct {
		name   string
		want   func() *filter.Filter
		mockFn func() *Telemetry
	}{
		{
			name: "Success",
			want: func() *filter.Filter {
				return nil
			},
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().Filter()
			assert.Equal(t, tt.want(), tr)
		})
	}
}

func TestTelemetry_Verbose(t *testing.T) {
	tests := []struct {
		name   string
		want   bool
		mockFn func() *Telemetry
	}{
		{
			name: "Success",
			want: false,
			mockFn: func() *Telemetry {
				return NewTelemetry()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn().Verbose()
			assert.Equal(t, tt.want, tr)
		})
	}
}
