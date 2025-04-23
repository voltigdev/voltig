/*
Package logger provides a logging interface for Voltig CLI.
*/
package logger

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var (
	// defaultLogger is the default logger instance.
	defaultLogger *log.Logger

	// LevelDebug is the debug log level.
	LevelDebug = log.DebugLevel
	// LevelInfo is the info log level.
	LevelInfo = log.InfoLevel
	// LevelWarn is the warn log level.
	LevelWarn = log.WarnLevel
	// LevelError is the error log level.
	LevelError = log.ErrorLevel
	// LevelFatal is the fatal log level.
	LevelFatal = log.FatalLevel
)

// Config holds logger configuration
type Config struct {
	Level      log.Level
	TimeFormat string
	Output     io.Writer
	Prefix     string
	ShowCaller bool
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      log.InfoLevel,
		TimeFormat: time.Kitchen,
		Output:     os.Stderr,
		Prefix:     "voltig",
		ShowCaller: false,
	}
}

// init initializes the default logger
func init() {
	cfg := DefaultConfig()
	defaultLogger = log.NewWithOptions(cfg.Output, log.Options{
		Level:           cfg.Level,
		ReportTimestamp: true,
		TimeFormat:      cfg.TimeFormat,
		Prefix:          cfg.Prefix,
		ReportCaller:    cfg.ShowCaller,
	})
}

// Configure sets up the logger with custom configuration
func Configure(cfg Config) {
	defaultLogger.SetLevel(cfg.Level)
	defaultLogger.SetOutput(cfg.Output)
	defaultLogger.SetPrefix(cfg.Prefix)
	defaultLogger.SetReportCaller(cfg.ShowCaller)
	defaultLogger.SetTimeFormat(cfg.TimeFormat)
}

// SetLevel sets the logging level
func SetLevel(level log.Level) {
	defaultLogger.SetLevel(level)
}

// Debug logs a debug message
func Debug(msg string, keyvals ...interface{}) {
	defaultLogger.Debug(msg, keyvals...)
}

// Info logs an info message
func Info(msg string, keyvals ...interface{}) {
	defaultLogger.Info(msg, keyvals...)
}

// Warn logs a warning message
func Warn(msg string, keyvals ...interface{}) {
	defaultLogger.Warn(msg, keyvals...)
}

// Error logs an error message
func Error(msg string, keyvals ...interface{}) {
	defaultLogger.Error(msg, keyvals...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, keyvals ...interface{}) {
	defaultLogger.Fatal(msg, keyvals...)
}

// WithPrefix returns a new logger with the given prefix
func WithPrefix(prefix string) *log.Logger {
	return defaultLogger.WithPrefix(prefix)
}

// WithValues returns a new logger with the given key-value pairs
func WithValues(keyvals ...interface{}) *log.Logger {
	return defaultLogger.With(keyvals...)
}
