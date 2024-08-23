package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewZapLogger(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		zapLogger, err := NewZapLogger()

		assert.NoError(t, err)
		assert.NotNil(t, zapLogger)
		assert.NotNil(t, zapLogger.logger)
	})
}

func TestZapLogger(t *testing.T) {
	tests := []struct {
		name     string
		method   func(l Logger, ctx context.Context, message string, err error, fields ...Field)
		message  string
		err      error
		fields   []Field
		expected map[string]interface{}
	}{
		{
			name: "Debug level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			message:  "Debug message",
			fields:   []Field{String("key1", "value1")},
			expected: map[string]interface{}{"key1": "value1"},
		},
		{
			name: "Info level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Info(ctx, message, fields...)
			},
			message:  "Info message",
			fields:   []Field{Int64("key2", 123)},
			expected: map[string]interface{}{"key2": int64(123)},
		},
		{
			name: "Warn level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Warn(ctx, message, fields...)
			},
			message:  "Warn message",
			fields:   []Field{Bool("key3", true)},
			expected: map[string]interface{}{"key3": true},
		},
		{
			name: "Error level with error",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Error(ctx, message, err, fields...)
			},
			message:  "Error message",
			err:      errors.New("an error occurred"),
			fields:   []Field{Float64("key4", 3.14)},
			expected: map[string]interface{}{"key4": 3.14, "error": "an error occurred"},
		},
		{
			name: "WithFields",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.WithFields(fields...).Debug(ctx, message)
			},
			message:  "Message with additional fields",
			fields:   []Field{String("additional", "info")},
			expected: map[string]interface{}{"additional": "info"},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			core, observedLogs := observer.New(zap.DebugLevel)
			logger := &ZapLogger{logger: zap.New(core)}

			// Clear previous entries
			observedLogs.TakeAll()

			// Run the logger method
			tt.method(logger, context.Background(), tt.message, tt.err, tt.fields...)

			// Get the observed logs
			logs := observedLogs.All()

			// Assert the log entry
			if assert.Equal(t, 1, len(logs)) {
				entry := logs[0]
				for key, value := range tt.expected {
					assert.Equal(t, value, entry.ContextMap()[key])
				}
				assert.Equal(t, tt.message, entry.Message)
			}
		})
	}
}
