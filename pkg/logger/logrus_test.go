package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger(t *testing.T) {
	tests := []struct {
		name string
		want *LogrusLogger
	}{
		{name: "Success", want: &LogrusLogger{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogrusLogger()
			assert.IsType(t, &LogrusLogger{}, got)
		})
	}
}

func TestLogrusLogger(t *testing.T) {
	tests := []struct {
		name     string
		method   func(l Logger, ctx context.Context, message string, err error, fields ...Field)
		message  string
		err      error
		fields   []Field
		expected logrus.Fields
	}{
		{
			name: "Debug level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			message:  "Debug message",
			fields:   []Field{String("key1", "value1")},
			expected: logrus.Fields{"key1": "value1"},
		},
		{
			name: "Info level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Info(ctx, message, fields...)
			},
			message:  "Info message",
			fields:   []Field{Int("key2", 123)},
			expected: logrus.Fields{"key2": 123},
		},
		{
			name: "Warn level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Warn(ctx, message, fields...)
			},
			message:  "Warn message",
			fields:   []Field{Bool("key3", true)},
			expected: logrus.Fields{"key3": true},
		},
		{
			name: "Error level with error",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Error(ctx, message, err, fields...)
			},
			message:  "Error message",
			err:      errors.New("an error occurred"),
			fields:   []Field{Float64("key4", 3.14)},
			expected: logrus.Fields{"key4": 3.14, "error": "an error occurred"},
		},
		{
			name: "WithFields",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.WithFields(fields...).Debug(ctx, message)
			},
			message:  "Message with additional fields",
			fields:   []Field{String("additional", "info")},
			expected: logrus.Fields{"additional": "info"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			baseLogger, hook := logrustest.NewNullLogger()
			baseLogger.SetLevel(logrus.DebugLevel)
			logger := &LogrusLogger{
				entry: logrus.NewEntry(baseLogger),
			}

			tt.method(logger, context.Background(), tt.message, tt.err, tt.fields...)

			if assert.Equal(t, 1, len(hook.Entries)) {
				entry := hook.LastEntry()
				assert.Equal(t, tt.expected, entry.Data)
				assert.Equal(t, tt.message, entry.Message)
			}

			hook.Reset()
		})
	}
}
