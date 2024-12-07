package logger

import (
	"context"
	"errors"
	"os"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/requestid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	correlationIDLabel = "_cID"
	traceIDLabel       = "trace-id"
	spanIDLabel        = "span-id"
)

type ZapLogger struct {
	logger *zap.Logger
	luja   *lumberjack.Logger
}

func NewZapLogger(serviceName string, lvl Level) *ZapLogger {
	var level zapcore.Level
	switch lvl {
	case DebugLevel:
		level = zap.DebugLevel
	case InfoLevel:
		level = zap.InfoLevel
	case WarnLevel:
		level = zap.WarnLevel
	case ErrorLevel:
		level = zap.ErrorLevel
	}

	z := zap.NewProductionConfig()
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	z.EncoderConfig.LevelKey = "severity"

	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewJSONEncoder(z.EncoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(level),
		),
	}

	var luja *lumberjack.Logger

	if serviceName != "" {
		luja = &lumberjack.Logger{
			Filename:   "logs/" + serviceName, // Log file location
			MaxSize:    10,                    // Max megabytes before log is rotated
			MaxBackups: 3,                     // Max number of old log files to keep
			MaxAge:     7,                     // Max number of days to retain old files
			Compress:   true,
		}

		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(z.EncoderConfig),
			zapcore.AddSync(luja),
			zap.NewAtomicLevelAt(level),
		))
	}

	return &ZapLogger{
		logger: zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1)),
		luja:   luja,
	}
}

func (z *ZapLogger) Debug(ctx context.Context, message string, fields ...Field) {
	zf := z.convertFields(fields)
	zf = append(zf, z.withTelemetry(ctx)...)
	z.logger.Debug(message, zf...)
}

func (z *ZapLogger) Info(ctx context.Context, message string, fields ...Field) {
	zf := z.convertFields(fields)
	zf = append(zf, z.withTelemetry(ctx)...)
	z.logger.Info(message, zf...)
}

func (z *ZapLogger) Warn(ctx context.Context, message string, fields ...Field) {
	zf := z.convertFields(fields)
	zf = append(zf, z.withTelemetry(ctx)...)
	z.logger.Warn(message, zf...)
}

func (z *ZapLogger) Error(ctx context.Context, message string, err error, fields ...Field) {
	zf := z.convertFields(fields)
	zf = append(zf, z.withTelemetry(ctx)...)
	zf = append(zf, zap.Error(err))
	z.logger.Error(message, zf...)
}

func (z *ZapLogger) WithFields(fields ...Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(z.convertFields(fields)...),
	}
}

func (z *ZapLogger) Close() error {
	var err error
	if z.luja != nil {
		err = errors.Join(z.luja.Close())
	}

	return errors.Join(err, z.logger.Sync())
}

func (z *ZapLogger) withTelemetry(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	spanCtx := trace.SpanFromContext(ctx).SpanContext()

	sid := spanCtx.SpanID().String()
	if sid == "0000000000000000" {
		sid = ""
	}

	tid := spanCtx.TraceID().String()
	if tid == "00000000000000000000000000000000" {
		tid = ""
	}

	return []zap.Field{
		zap.String(correlationIDLabel, requestid.Get(ctx)),
		zap.String(spanIDLabel, sid),
		zap.String(traceIDLabel, tid),
	}
}

func (z *ZapLogger) convertFields(fields []Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.key, field.value)
	}

	return zapFields
}
