package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(level Level) (*ZapLogger, error) {
	var lvl zapcore.Level
	switch level {
	case DebugLevel:
		lvl = zap.DebugLevel
	case InfoLevel:
		lvl = zap.InfoLevel
	case WarnLevel:
		lvl = zap.WarnLevel
	case ErrorLevel:
		lvl = zap.ErrorLevel
	default:
		lvl = zap.InfoLevel
	}

	z := zap.NewProductionConfig()
	z.DisableCaller = true
	z.Level = zap.NewAtomicLevelAt(lvl)
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	z.EncoderConfig.LevelKey = "severity"

	logger, err := z.Build()
	if err != nil {
		return nil, err
	}

	logger = zap.New(logger.Core(), zap.AddCaller(), zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger}, nil
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

func (z *ZapLogger) Close() error {
	return z.logger.Sync()
}

func (z *ZapLogger) convertFields(fields []Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}

	return zapFields
}
