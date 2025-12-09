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

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...any) {
	if l.level <= LevelDebug {
		message := fmt.Sprintf(format, args...)
		l.log("DEBUG", message, false)
	}
}

func (l *Logger) Highlight(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.formatter.FormatHighlight(message)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.log("INFO", message, false)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	if l.level <= LevelWarn {
		message := fmt.Sprintf(format, args...)
		l.log("WARN", message, false)
	}
}

// Error logs an error message to stderr
func (l *Logger) Error(format string, args ...any) {
	if l.level <= LevelError {
		message := fmt.Sprintf(format, args...)
		l.log("ERROR", message, true)
	}
}

// Success logs a success message
func (l *Logger) Success(format string, args ...any) {
	if l.level <= LevelInfo {
		message := fmt.Sprintf(format, args...)
		l.log("SUCCESS", message, false)
	}
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.log("FATAL", message, true)
	os.Exit(1)
}

// log is the internal logging method
func (l *Logger) log(level, message string, isError bool) {
	tag := level
	if l.prefix != "" {
		tag = fmt.Sprintf("%s:%s", l.prefix, level)
	}

	if isError {
		l.formatter.FormatError(tag, message)
	} else {
		l.formatter.FormatSuccess(tag, message)
	}
}

// Printf logs a formatted message at info level
func (l *Logger) Printf(format string, args ...any) {
	l.Info(format, args...)
}

// Println logs a message at info level
func (l *Logger) Println(args ...any) {
	message := fmt.Sprint(args...)
	l.Info("%s", message)
}

// Output logs a message with a custom tag/prefix
func (l *Logger) Output(tag, message string, isError bool) {
	if isError {
		l.formatter.FormatError(tag, message)
	} else {
		l.formatter.FormatSuccess(tag, message)
	}
}

// FormatTable formats a table of data
func (l *Logger) FormatTable(headers []string, rows [][]string) string {
	return l.formatter.FormatTable(headers, rows)
}
