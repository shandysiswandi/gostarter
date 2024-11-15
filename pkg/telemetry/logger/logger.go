package logger

import "context"

// Logger Interface. All methods SHOULD be safe for concurrent use.
type Logger interface {
	// Debug logs a message at the Debug level with optional structured fields.
	Debug(ctx context.Context, message string, fields ...Field)

	// Info logs a message at the Info level with optional structured fields.
	Info(ctx context.Context, message string, fields ...Field)

	// Warn logs a message at the Warn level with optional structured fields.
	Warn(ctx context.Context, message string, fields ...Field)

	// Error logs a message at the Error level with an error and optional structured fields.
	Error(ctx context.Context, message string, err error, fields ...Field)

	// WithFields returns a new Logger instance with the specified fields added to each log message.
	WithFields(fields ...Field) Logger
}

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

type Field struct {
	key   string
	value any
}

// KeyVal creates a Field for any value.
func KeyVal(key string, value any) Field {
	return Field{key: key, value: value}
}
