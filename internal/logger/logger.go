package logger

import (
	"fmt"
	"os"
	"time"
)

// Level represents log level.
type Level int

const (
	// LevelDebug represents debug log level.
	LevelDebug Level = iota
	// LevelInfo represents info log level.
	LevelInfo
	// LevelWarn represents warning log level.
	LevelWarn
	// LevelError represents error log level.
	LevelError
)

// Logger provides structured logging functionality.
type Logger struct {
	level Level
	debug bool
}

// New creates a new logger instance.
func New(debug bool) *Logger {
	level := LevelInfo
	if debug {
		level = LevelDebug
	}
	return &Logger{level: level, debug: debug}
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", format, args...)
	}
}

// Info logs an info message.
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", format, args...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", format, args...)
	}
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.log("ERROR", format, args...)
	}
}

// log writes a log message to stderr.
func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, level, msg)
}


