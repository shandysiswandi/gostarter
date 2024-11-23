package telemetry

import (
	"errors"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNewTelemetry(t *testing.T) {
	type args struct {
		opts []func(*Telemetry)
	}
	tests := []struct {
		name string
		args args
		want *Telemetry
	}{
		{
			name: "Success",
			args: args{
				opts: []func(*Telemetry){
					func(t *Telemetry) {},
				},
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
			wantErr: errors.New("sync /dev/stderr: invalid argument"),
			mockFn: func() *Telemetry {
				return NewTelemetry(WithZapLogger(logger.InfoLevel))
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
				return noop.NewTracerProvider().Tracer("telemetry")
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
			name: "SuccessNoopTracer",
			want: func() trace.TracerProvider {
				return noop.NewTracerProvider()
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
		want   func() Collector
		mockFn func() *Telemetry
	}{
		{
			name: "SuccessNoopTracer",
			want: func() Collector {
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
			tr := tt.mockFn().TracerCollector()
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
			tr := tt.mockFn().Meter()
			assert.Equal(t, tt.want(), tr)
		})
	}
}
