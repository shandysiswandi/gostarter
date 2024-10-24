package logger

import (
	"context"
)

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (z *NoopLogger) Debug(context.Context, string, ...Field) {
	// do nothing
}

func (z *NoopLogger) Info(context.Context, string, ...Field) {
	// do nothing
}

func (z *NoopLogger) Warn(context.Context, string, ...Field) {
	// do nothing
}

func (z *NoopLogger) Error(context.Context, string, error, ...Field) {
	// do nothing
}

func (z *NoopLogger) WithFields(...Field) Logger {
	return &NoopLogger{}
}
