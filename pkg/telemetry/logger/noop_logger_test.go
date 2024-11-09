package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNoopLogger(t *testing.T) {
	tests := []struct {
		name string
		want *NoopLogger
	}{
		{
			name: "Success",
			want: &NoopLogger{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewNoopLogger()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoopLogger_Debug(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 string
		in2 []Field
	}
	tests := []struct {
		name string
		z    *NoopLogger
		args args
	}{
		{
			name: "Success",
			z:    &NoopLogger{},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := &NoopLogger{}
			z.Debug(tt.args.in0, tt.args.in1, tt.args.in2...)
		})
	}
}

func TestNoopLogger_Info(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 string
		in2 []Field
	}
	tests := []struct {
		name string
		z    *NoopLogger
		args args
	}{
		{
			name: "Success",
			z:    &NoopLogger{},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := &NoopLogger{}
			z.Info(tt.args.in0, tt.args.in1, tt.args.in2...)
		})
	}
}

func TestNoopLogger_Warn(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 string
		in2 []Field
	}
	tests := []struct {
		name string
		z    *NoopLogger
		args args
	}{
		{
			name: "Success",
			z:    &NoopLogger{},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := &NoopLogger{}
			z.Warn(tt.args.in0, tt.args.in1, tt.args.in2...)
		})
	}
}

func TestNoopLogger_Error(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 string
		in2 error
		in3 []Field
	}
	tests := []struct {
		name string
		z    *NoopLogger
		args args
	}{
		{
			name: "Success",
			z:    &NoopLogger{},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := &NoopLogger{}
			z.Error(tt.args.in0, tt.args.in1, tt.args.in2, tt.args.in3...)
		})
	}
}

func TestNoopLogger_WithFields(t *testing.T) {
	type args struct {
		in0 []Field
	}
	tests := []struct {
		name string
		z    *NoopLogger
		args args
		want Logger
	}{
		{
			name: "Success",
			z:    &NoopLogger{},
			args: args{},
			want: &NoopLogger{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := &NoopLogger{}
			got := z.WithFields(tt.args.in0...)
			assert.Equal(t, tt.want, got)
		})
	}
}
