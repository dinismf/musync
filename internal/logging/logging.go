package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogLevel represents the level of logging
type LogLevel int

// Log levels
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger is a wrapper around the standard log package
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       LogLevel
}

// Global logger instance
var logger *Logger

// Init initializes the logger with the given log level
func Init(level LogLevel, output io.Writer) {
	if output == nil {
		output = os.Stdout
	}

	// Create loggers for each level
	debugLogger := log.New(output, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger := log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger := log.New(output, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger := log.New(output, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Create our logger wrapper
	logger = &Logger{
		debugLogger: debugLogger,
		infoLogger:  infoLogger,
		warnLogger:  warnLogger,
		errorLogger: errorLogger,
		fatalLogger: fatalLogger,
		level:       level,
	}
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if logger == nil {
		// If the logger hasn't been initialized, create a default one
		Init(InfoLevel, os.Stdout)
	}
	return logger
}

// WithContext returns a logger with the given context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract any context values that should be added to the logger
	// For now, we're just returning the same logger
	return l
}

// Debug logs a message at debug level
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info logs a message at info level
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn logs a message at warn level
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error logs a message at error level
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		l.fatalLogger.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method, path, clientIP string, statusCode int, latency time.Duration) {
	l.Info("[%s] %s %s %d %s", method, path, clientIP, statusCode, latency)
}