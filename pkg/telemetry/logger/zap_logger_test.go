package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestZapWithVerbose(t *testing.T) {
	type args struct {
		isVerbose bool
	}
	tests := []struct {
		name   string
		args   args
		want   bool
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{isVerbose: true},
			want: true,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithVerbose(a.isVerbose)

				z, _ := NewZapLogger(o)

				return z
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.mockFn(tt.args)

			assert.Equal(t, tt.want, z.option.isVerbose)
		})
	}
}

func TestZapWithFilteredKeys(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		args   args
		want   int
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Success",
			args: args{keys: []string{"header", "auth"}},
			want: 2,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithFilteredKeys(a.keys)

				z, _ := NewZapLogger(o)

				return z
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.mockFn(tt.args)

			assert.Len(t, z.option.filteredKeys, tt.want)
		})
	}
}

func TestZapWithLevel(t *testing.T) {
	type args struct {
		lvl Level
	}
	tests := []struct {
		name   string
		args   args
		want   zapcore.Level
		mockFn func(a args) *ZapLogger
	}{
		{
			name: "Debug",
			args: args{lvl: DebugLevel},
			want: zapcore.DebugLevel,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithLevel(a.lvl)

				z, _ := NewZapLogger(o)

				return z
			},
		},
		{
			name: "Info",
			args: args{lvl: InfoLevel},
			want: zapcore.InfoLevel,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithLevel(a.lvl)

				z, _ := NewZapLogger(o)

				return z
			},
		},
		{
			name: "Warn",
			args: args{lvl: WarnLevel},
			want: zapcore.WarnLevel,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithLevel(a.lvl)

				z, _ := NewZapLogger(o)

				return z
			},
		},
		{
			name: "Error",
			args: args{lvl: ErrorLevel},
			want: zapcore.ErrorLevel,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithLevel(a.lvl)

				z, _ := NewZapLogger(o)

				return z
			},
		},
		{
			name: "Default",
			args: args{lvl: DebugLevel - 1},
			want: zapcore.InfoLevel,
			mockFn: func(a args) *ZapLogger {
				o := ZapWithLevel(a.lvl)

				z, _ := NewZapLogger(o)

				return z
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			z := tt.mockFn(tt.args)

			assert.Equal(t, tt.want, z.option.level)
		})
	}
}

func TestNewZapLogger(t *testing.T) {
	type args struct {
		opts []func(*ZapOption)
	}
	tests := []struct {
		name    string
		args    args
		want    *ZapLogger
		wantErr error
	}{
		{
			name: "Success",
			args: args{},
			want: &ZapLogger{
				option: &ZapOption{
					level:        zapcore.InfoLevel,
					filteredKeys: make([]string, 0),
					isVerbose:    true,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewZapLogger(tt.args.opts...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want.option, got.option)
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
			name: "Success",
			args: args{
				ctx:     context.Background(),
				message: "debug",
				fields:  []Field{},
			},
			mockFn: func(a args) *ZapLogger {
				z, _ := NewZapLogger()
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
				z, _ := NewZapLogger()
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
				z, _ := NewZapLogger()
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
				z, _ := NewZapLogger(ZapWithFilteredKeys([]string{"b"}))
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
	z, _ := NewZapLogger()

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
				z, _ := NewZapLogger()

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
