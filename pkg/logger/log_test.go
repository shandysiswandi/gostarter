package logger

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"
)

func TestLogLogger(t *testing.T) {
	tests := []struct {
		name     string
		method   func(l Logger, ctx context.Context, message string, err error, fields ...Field)
		level    string
		message  string
		err      error
		fields   []Field
		expected string
	}{
		{
			name: "Debug level String",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{String("key1", "value1")},
			expected: "DEBUG: Debug message | key1=value1 \n",
		},
		{
			name: "Debug level Bool",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Bool("key1", false)},
			expected: "DEBUG: Debug message | key1=false \n",
		},
		{
			name: "Debug level Int",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Int("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Int8",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Int8("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Int16",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Int16("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Int32",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Int32("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Int64",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Int64("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Uint",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Uint("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Uint8",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Uint8("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Uint16",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Uint16("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Uint32",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Uint32("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Uint64",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Uint64("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Float32",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Float32("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Float64",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Float64("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Debug level Bytes",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Bytes("key1", []byte(`{"key":"value"}`))},
			expected: "DEBUG: Debug message | key1={\"key\":\"value\"} \n",
		},
		{
			name: "Debug level Any",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Debug(ctx, message, fields...)
			},
			level:    "DEBUG",
			message:  "Debug message",
			fields:   []Field{Any("key1", 1)},
			expected: "DEBUG: Debug message | key1=1 \n",
		},
		{
			name: "Info level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Info(ctx, message, fields...)
			},
			level:    "INFO",
			message:  "Info message",
			fields:   []Field{Int("key2", 123)},
			expected: "INFO: Info message | key2=123 \n",
		},
		{
			name: "Warn level",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Warn(ctx, message, fields...)
			},
			level:    "WARN",
			message:  "Warn message",
			fields:   []Field{Bool("key3", true)},
			expected: "WARN: Warn message | key3=true \n",
		},
		{
			name: "Error level with error",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.Error(ctx, message, err, fields...)
			},
			level:    "ERROR",
			message:  "Error message",
			err:      errors.New("an error occurred"),
			fields:   []Field{Float64("key4", 3.14)},
			expected: "ERROR: Error message ERR: an error occurred | key4=3.14 \n",
		},
		{
			name: "WithFields",
			method: func(l Logger, ctx context.Context, message string, err error, fields ...Field) {
				l.WithFields(fields...).Debug(ctx, message)
			},
			level:    "DEBUG",
			message:  "Message with additional fields",
			fields:   []Field{String("additional", "info")},
			expected: "DEBUG: Message with additional fields | additional=info \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			logger := &LogLogger{logger: log.New(&buf, "", 0)}

			tt.method(logger, context.Background(), tt.message, tt.err, tt.fields...)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
