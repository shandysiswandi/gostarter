// Package logger provides an abstraction for structured logging.
//
// It defines a `Logger` interface for logging messages at different levels with optional fields.
// Additionally, it provides functions to create `Field` instances with various types of values.
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
)

// LogLogger is a concrete implementation of the Logger interface using the standard library's log package.
type LogLogger struct {
	logger *log.Logger
	mu     sync.Mutex
	fields []Field
}

// NewStdLogger creates a new instance of LogLogger with the standard log.Logger.
func NewStdLogger() *LogLogger {
	return &LogLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// log is a helper method to format and output log messages with a specific log level.
func (l *LogLogger) log(_ context.Context, level string, message string, err error, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Combine base fields with the additional fields
	allFields := append([]Field(nil), l.fields...)
	allFields = append(allFields, fields...)

	// Build the log message
	logMessage := level + ": " + message
	if err != nil {
		logMessage = logMessage + " ERR: " + err.Error()
	}
	if len(allFields) > 0 {
		logMessage += " | "
		for _, field := range allFields {
			logMessage += field.Key + "=" + formatValue(field.Value) + " "
		}
	}

	// Output the log message
	//nolint:errcheck // no error happened here, maybe
	_ = l.logger.Output(3, logMessage)
}

// Debug logs a message at the Debug level with optional fields.
func (l *LogLogger) Debug(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, "DEBUG", message, nil, fields...)
}

// Info logs a message at the Info level with optional fields.
func (l *LogLogger) Info(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, "INFO", message, nil, fields...)
}

// Warn logs a message at the Warn level with optional fields.
func (l *LogLogger) Warn(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, "WARN", message, nil, fields...)
}

// Error logs a message at the Error level with an error and optional fields.
func (l *LogLogger) Error(ctx context.Context, message string, err error, fields ...Field) {
	l.log(ctx, "ERROR", message, err, fields...)
}

// WithFields returns a new LogLogger instance with the specified fields added to each log message.
func (l *LogLogger) WithFields(fields ...Field) Logger {
	return &LogLogger{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

// formatValue formats a value into a string suitable for logging.
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
