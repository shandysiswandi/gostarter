package logger

import (
	"context"
	"slices"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/telemetry/requestid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	correlationIDLabel = "_cID"
	traceIDLabel       = "trace-id"
	spanIDLabel        = "span-id"
)

type ZapOption struct {
	level        zapcore.Level
	filteredKeys []string
	isVerbose    bool
}

type ZapLogger struct {
	logger *zap.Logger
	option *ZapOption
}

func ZapWithVerbose(isVerbose bool) func(*ZapOption) {
	return func(zo *ZapOption) {
		zo.isVerbose = isVerbose
	}
}

func ZapWithFilteredKeys(keys []string) func(*ZapOption) {
	return func(zo *ZapOption) {
		zo.filteredKeys = append(zo.filteredKeys, keys...)
	}
}

func ZapWithLevel(lvl Level) func(*ZapOption) {
	return func(zo *ZapOption) {
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
		default:
			level = zap.InfoLevel
		}

		zo.level = level
	}
}

func NewZapLogger(opts ...func(*ZapOption)) (*ZapLogger, error) {
	zopt := &ZapOption{level: zap.InfoLevel, isVerbose: true, filteredKeys: make([]string, 0)}

	for _, opt := range opts {
		opt(zopt)
	}

	z := zap.NewProductionConfig()
	z.DisableCaller = true
	z.Level = zap.NewAtomicLevelAt(zopt.level)
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	z.EncoderConfig.LevelKey = "severity"

	logger, err := z.Build()
	if err != nil {
		return nil, err
	}

	logger = zap.New(logger.Core(), zap.AddCaller(), zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger, option: zopt}, nil
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
		option: z.option,
	}
}

func (z *ZapLogger) Close() error {
	return z.logger.Sync()
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
		if ok := slices.Contains(z.option.filteredKeys, strings.ToLower(field.key)); ok {
			zapFields[i] = zap.String(field.key, "***")

			continue
		}

		zapFields[i] = zap.Any(field.key, field.value)
	}

	return zapFields
}
