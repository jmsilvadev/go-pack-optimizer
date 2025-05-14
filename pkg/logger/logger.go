package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field is an alias for zapcore.Field, used to add structured data to logs.
type Field = zapcore.Field

// Logger defines the interface for a structured logger.
// It supports common log levels with optional structured fields.
type Logger interface {
	Debug(msg string, fields ...Field) // Logs a debug-level message
	Info(msg string, fields ...Field)  // Logs an info-level message
	Warn(msg string, fields ...Field)  // Logs a warning-level message
	Error(msg string, fields ...Field) // Logs an error-level message
}

// Supported log levels
const (
	LEVEL_DEBUG = zapcore.DebugLevel
	LEVEL_INFO  = zapcore.InfoLevel
	LEVEL_WARN  = zapcore.WarnLevel
	LEVEL_ERROR = zapcore.ErrorLevel
)

// New initializes and returns a new zap-based Logger implementation.
// It writes structured logs in JSON format to stdout with the specified log level.
func New(level zapcore.Level) Logger {
	consoleErrors := zapcore.Lock(os.Stdout)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.MillisDurationEncoder
	config.EncodeName = zapcore.FullNameEncoder

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(level)

	consoleEncoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, logLevel),
	)

	caller := zap.AddCaller()
	dev := zap.Development()

	log := zap.New(core, caller, dev)

	return log
}
