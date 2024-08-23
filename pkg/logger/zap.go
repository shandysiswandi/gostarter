// Package logger provides an abstraction for structured logging.
//
// It defines a `Logger` interface for logging messages at different levels with optional fields.
// Additionally, it provides functions to create `Field` instances with various types of values.
package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	baseLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: baseLogger}, nil
}

func (z *ZapLogger) Debug(_ context.Context, message string, fields ...Field) {
	z.logger.Debug(message, z.convertFields(fields)...)
}

func (z *ZapLogger) Info(_ context.Context, message string, fields ...Field) {
	z.logger.Info(message, z.convertFields(fields)...)
}

func (z *ZapLogger) Warn(_ context.Context, message string, fields ...Field) {
	z.logger.Warn(message, z.convertFields(fields)...)
}

func (z *ZapLogger) Error(_ context.Context, message string, err error, fields ...Field) {
	zf := z.convertFields(fields)
	zf = append(zf, zap.Error(err))
	z.logger.Error(message, zf...)
}

func (z *ZapLogger) WithFields(fields ...Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(z.convertFields(fields)...),
	}
}

func (z *ZapLogger) convertFields(fields []Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}

	return zapFields
}
