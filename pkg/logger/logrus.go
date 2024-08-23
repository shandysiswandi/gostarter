// Package logger provides an abstraction for structured logging.
//
// It defines a `Logger` interface for logging messages at different levels with optional fields.
// Additionally, it provides functions to create `Field` instances with various types of values.
package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

// LogrusLogger is a concrete implementation of the Logger interface using logrus.
// It provides structured logging with support for various log levels and fields.
type LogrusLogger struct {
	entry *logrus.Entry
}

// NewLogrusLogger creates a new instance of LogrusLogger.
// It initializes the underlying logrus.Logger with a text formatter that includes full timestamps.
func NewLogrusLogger() *LogrusLogger {
	// Create a new logrus logger instance
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return &LogrusLogger{
		entry: logrus.NewEntry(baseLogger),
	}
}

// Debug logs a message at the Debug level with optional fields.
// It uses the context to enrich log entries and converts the custom fields to logrus.Fields.
func (l *LogrusLogger) Debug(ctx context.Context, message string, fields ...Field) {
	l.entry.WithContext(ctx).WithFields(l.convertFields(fields)).Debug(message)
}

// Info logs a message at the Info level with optional fields.
// It uses the context to enrich log entries and converts the custom fields to logrus.Fields.
func (l *LogrusLogger) Info(ctx context.Context, message string, fields ...Field) {
	l.entry.WithContext(ctx).WithFields(l.convertFields(fields)).Info(message)
}

// Warn logs a message at the Warn level with optional fields.
// It uses the context to enrich log entries and converts the custom fields to logrus.Fields.
func (l *LogrusLogger) Warn(ctx context.Context, message string, fields ...Field) {
	l.entry.WithContext(ctx).WithFields(l.convertFields(fields)).Warn(message)
}

// Error logs a message at the Error level with an error and optional fields.
// It uses the context to enrich log entries and converts the custom fields to logrus.Fields.
func (l *LogrusLogger) Error(ctx context.Context, message string, err error, fields ...Field) {
	lf := l.convertFields(fields)
	if err != nil {
		lf["error"] = err.Error()
	}
	l.entry.WithContext(ctx).WithFields(lf).Error(message)
}

// WithFields returns a new LogrusLogger instance with the specified fields added to each log message.
// The new instance inherits the base logger's configuration and includes additional fields.
func (l *LogrusLogger) WithFields(fields ...Field) Logger {
	return &LogrusLogger{
		entry: l.entry.WithFields(l.convertFields(fields)),
	}
}

// convertFields converts the custom Field type to logrus.Fields.
// This method maps each Field's key and value to the corresponding logrus.Fields format.
func (l *LogrusLogger) convertFields(fields []Field) logrus.Fields {
	logrusFields := logrus.Fields{}
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}

	return logrusFields
}
