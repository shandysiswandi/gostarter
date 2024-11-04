package logger

import "context"

// Logger Interface. All methods SHOULD be safe for concurrent use.
type Logger interface {
	// Debug logs a message at the Debug level with optional structured fields.
	//
	// Parameters:
	// - ctx: The context for the logging operation.
	// - message: The message to log.
	// - fields: Optional fields to include with the log message.
	//
	// Example:
	// logger.Debug(ctx, "This is a debug message", logger.String("key", "value"))
	Debug(ctx context.Context, message string, fields ...Field)

	// Info logs a message at the Info level with optional structured fields.
	//
	// Parameters:
	// - ctx: The context for the logging operation.
	// - message: The message to log.
	// - fields: Optional fields to include with the log message.
	//
	// Example:
	// logger.Info(ctx, "This is an info message", logger.Int("count", 42))
	Info(ctx context.Context, message string, fields ...Field)

	// Warn logs a message at the Warn level with optional structured fields.
	//
	// Parameters:
	// - ctx: The context for the logging operation.
	// - message: The message to log.
	// - fields: Optional fields to include with the log message.
	//
	// Example:
	// logger.Warn(ctx, "This is a warning message", logger.Bool("flag", true))
	Warn(ctx context.Context, message string, fields ...Field)

	// Error logs a message at the Error level with an error and optional structured fields.
	//
	// Parameters:
	// - ctx: The context for the logging operation.
	// - message: The message to log.
	// - err: The error to log along with the message.
	// - fields: Optional fields to include with the log message.
	//
	// Example:
	// logger.Error(ctx, "This is an error message", errors.New("an error occurred"), logger.Float64("value", 3.14))
	Error(ctx context.Context, message string, err error, fields ...Field)

	// WithFields returns a new Logger instance with the specified fields added to each log message.
	//
	// Parameters:
	// - fields: The fields to include with each log message.
	//
	// Returns:
	// - A new Logger instance that includes the specified fields in every log message.
	//
	// Example:
	// newLogger := logger.WithFields(logger.String("session", "12345"))
	// newLogger.Info(ctx, "User logged in")
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
