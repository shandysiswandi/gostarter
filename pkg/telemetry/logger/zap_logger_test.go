package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZapLogger(t *testing.T) {
	type args struct {
		filename string
		lvl      Level
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewZapLogger(tt.args.filename, tt.args.lvl)
			assert.NotNil(t, got)
		})
	}
}

func TestZapLogger_Debug(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
		fields  []Field
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "SuccessButContextNil",
			args: args{
				ctx:     nil,
				message: "debug",
				fields:  []Field{},
			},
			mockFn: func(a args) *ZapLogger {
				return NewZapLogger("", DebugLevel)
			},
		},
		{
			name: "Success",
			args: args{
				ctx:     context.Background(),
				message: "debug",
				fields:  []Field{},
			},
			mockFn: func(a args) *ZapLogger {
				return NewZapLogger("", DebugLevel)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn(tt.args)
			z.Debug(tt.args.ctx, tt.args.message, tt.args.fields...)
		})
	}
}

func TestZapLogger_Info(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
		fields  []Field
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{
				ctx:     context.Background(),
				message: "info",
				fields:  []Field{},
			},
			mockFn: func(a args) *ZapLogger {
				return NewZapLogger("", InfoLevel)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn(tt.args)
			z.Info(tt.args.ctx, tt.args.message, tt.args.fields...)
		})
	}
}

func TestZapLogger_Warn(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
		fields  []Field
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{
				ctx:     context.Background(),
				message: "warn",
				fields:  []Field{},
			},
			mockFn: func(a args) *ZapLogger {
				return NewZapLogger("", WarnLevel)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn(tt.args)
			z.Warn(tt.args.ctx, tt.args.message, tt.args.fields...)
		})
	}
}

func TestZapLogger_Error(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
		err     error
		fields  []Field
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{
				ctx:     context.Background(),
				message: "error",
				err:     assert.AnError,
				fields:  []Field{{key: "a", value: "a"}, {key: "b", value: "b"}},
			},
			mockFn: func(a args) *ZapLogger {
				return NewZapLogger("", ErrorLevel)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn(tt.args)
			z.Error(tt.args.ctx, tt.args.message, tt.args.err, tt.args.fields...)
		})
	}
}

func TestZapLogger_WithFields(t *testing.T) {
	z := NewZapLogger("", InfoLevel)

	type args struct {
		fields []Field
	}
	tests := []struct {
		name   string
		args   args
		want   Logger
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{fields: []Field{}},
			want: z,
			mockFn: func(a args) *ZapLogger {
				return z
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn(tt.args)
			got := z.WithFields(tt.args.fields...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestZapLogger_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *ZapLogger
	}{
		{
			name:    "Error",
			wantErr: errors.New("sync /dev/stdout: invalid argument"),
			mockFn: func() *ZapLogger {
				return NewZapLogger("", InfoLevel)
			},
		},
		{
			name:    "ErrorFile",
			wantErr: errors.New("sync /dev/stdout: invalid argument"),
			mockFn: func() *ZapLogger {
				return NewZapLogger("./.out", InfoLevel)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			z := tt.mockFn()
			err := z.Close()
			if err != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}
