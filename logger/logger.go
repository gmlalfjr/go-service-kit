package logger

import "fmt"

// Logger interface for logging
type Logger interface {
	Info(message string)
	Error(message string, err error)
}

// DefaultLogger is the default implementation of Logger interface
type DefaultLogger struct{}

// Info logs an information message
func (l *DefaultLogger) Info(message string) {
	fmt.Println("[INFO]", message)
}

// Error logs an error message
func (l *DefaultLogger) Error(message string, err error) {
	fmt.Println("[ERROR]", message, err)
}

// NewDefaultLogger creates a new instance of DefaultLogger
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}
