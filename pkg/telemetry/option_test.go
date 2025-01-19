package telemetry

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/filter"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/stretchr/testify/assert"
)

func TestWithServiceName(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{serviceName: "gostarter"},
			want: "gostarter",
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithServiceName(a.serviceName)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tel := tt.mockFn(tt.args)
			assert.Equal(t, tel.name, tt.want)
		})
	}
}

func TestWithLogFilter(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		args   args
		want   *filter.Filter
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{keys: []string{"token"}},
			want: filter.NewFilter(
				filter.WithHeaders("token"),
				filter.WithQueries("token"),
				filter.WithFields("token"),
			),
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithLogFilter(a.keys...)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tel := tt.mockFn(tt.args)
			assert.Equal(t, tel.filter, tt.want)
		})
	}
}

func TestWithZapLogger(t *testing.T) {
	type args struct {
		serviceName string
		level       logger.Level
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{
				serviceName: "gostarter",
				level:       logger.InfoLevel,
			},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithZapLogger("", logger.InfoLevel, false)(tel)

				return tel
			},
		},
		{
			name: "SuccessWithFile",
			args: args{
				serviceName: "gostarter",
				level:       logger.InfoLevel,
			},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithZapLogger("file", logger.InfoLevel, true)(tel)

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

func TestWithOTLPTracer(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Telemetry
	}{
		{
			name: "Success",
			args: args{address: "gostarter"},
			mockFn: func(a args) *Telemetry {
				tel := NewTelemetry()

				WithOTLP(a.address)(tel)

				return tel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tel := tt.mockFn(tt.args)
			defer tel.Close()

			assert.Len(t, tel.flushers, 3)
			assert.NotNil(t, tel.Tracer())
		})
	}
}
