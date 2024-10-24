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
	Key   string
	Value any
}

// String creates a Field for a string value.
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a Field for a bool value.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Int creates a Field for an int value.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int8 creates a Field for an int8 value.
func Int8(key string, value int8) Field {
	return Field{Key: key, Value: value}
}

// Int16 creates a Field for an int16 value.
func Int16(key string, value int16) Field {
	return Field{Key: key, Value: value}
}

// Int32 creates a Field for an int32 value.
func Int32(key string, value int32) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates a Field for an int64 value.
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Uint creates a Field for a uint value.
func Uint(key string, value uint) Field {
	return Field{Key: key, Value: value}
}

// Uint8 creates a Field for a uint8 value.
func Uint8(key string, value uint8) Field {
	return Field{Key: key, Value: value}
}

// Uint16 creates a Field for a uint16 value.
func Uint16(key string, value uint16) Field {
	return Field{Key: key, Value: value}
}

// Uint32 creates a Field for a uint32 value.
func Uint32(key string, value uint32) Field {
	return Field{Key: key, Value: value}
}

// Uint64 creates a Field for a uint64 value.
func Uint64(key string, value uint64) Field {
	return Field{Key: key, Value: value}
}

// Float32 creates a Field for a float32 value.
func Float32(key string, value float32) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a Field for a float64 value.
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bytes creates a Field for a []byte value.
func Bytes(key string, value []byte) Field {
	return Field{Key: key, Value: value}
}

// Any creates a Field for any value.
func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}
