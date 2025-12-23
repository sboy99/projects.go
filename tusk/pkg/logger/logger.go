package logger

import (
	"fmt"
	"os"

	"github.com/sboy99/projects.go/tusk/pkg/formatter"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger provides structured logging with formatting
type Logger struct {
	formatter *formatter.Formatter
	level     LogLevel
	prefix    string
}

// NewLogger creates a new logger instance
func NewLogger(indent string) *Logger {
	return &Logger{
		formatter: formatter.NewFormatter(indent),
		level:     LevelInfo,
		prefix:    "",
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetPrefix sets a prefix for all log messages
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Highlight(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.formatter.FormatHighlight(message)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...any) {
	if l.level <= LevelDebug {
		message := fmt.Sprintf(format, args...)
		l.log("DEBUG", message)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.log("INFO", message)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	if l.level <= LevelWarn {
		message := fmt.Sprintf(format, args...)
		l.log("WARN", message)
	}
}

// Error logs an error message to stderr
func (l *Logger) Error(format string, args ...any) {
	if l.level <= LevelError {
		message := fmt.Sprintf(format, args...)
		l.log("ERROR", message)
	}
}

// Success logs a success message
func (l *Logger) Success(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.log("SUCCESS", message)
	}
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.log("FATAL", message)
	os.Exit(1)
}

func (l *Logger) Plain(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.log("", message)
}

// log is the internal logging method
func (l *Logger) log(level, message string) {
	tag := level
	if l.prefix != "" {
		tag = fmt.Sprintf("%s:%s", l.prefix, level)
	}
	switch level {
	case "SUCCESS":
		l.formatter.FormatSuccess(tag, message)
	case "ERROR":
		l.formatter.FormatError(tag, message)
	case "WARN":
		l.formatter.FormatWarning(tag, message)
	case "INFO":
		l.formatter.FormatInfo(tag, message)
	case "DEBUG":
		l.formatter.FormatDebug(tag, message)
	case "FATAL":
		l.formatter.FormatError(tag, message)
	default:
		l.formatter.FormatPlain(tag, message)
	}
}

// Output logs a message with a custom tag/prefix
func (l *Logger) Output(tag, message string, isError bool) {
	if isError {
		l.formatter.FormatError(tag, message)
	} else {
		l.formatter.FormatPlain(tag, message)
	}
}

// FormatTable formats a table of data
func (l *Logger) FormatTable(headers []string, rows [][]string) {
	l.formatter.FormatTable(headers, rows)
}
