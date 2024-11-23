package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZapLogger(t *testing.T) {
	type args struct {
		lvl  Level
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    *ZapLogger
		wantErr error
	}{
		{
			name:    "Success",
			args:    args{},
			want:    &ZapLogger{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewZapLogger(tt.args.lvl, tt.args.keys...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want.filteredKeys, got.filteredKeys)
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
				z, _ := NewZapLogger(DebugLevel)
				return z
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
				z, _ := NewZapLogger(DebugLevel)
				return z
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
				z, _ := NewZapLogger(InfoLevel)
				return z
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
				z, _ := NewZapLogger(WarnLevel)
				return z
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
				z, _ := NewZapLogger(ErrorLevel, "b")
				return z
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
	z, _ := NewZapLogger(InfoLevel)

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
			wantErr: errors.New("sync /dev/stderr: invalid argument"),
			mockFn: func() *ZapLogger {
				z, _ := NewZapLogger(InfoLevel)

				return z
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
