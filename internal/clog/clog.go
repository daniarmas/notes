package clog

import (
	"log/slog"
	"os"
)

// Clog is a custom logger that wraps the slog package
type Clog struct {
	logger *slog.Logger
}

// NewClog creates a new Clog instance
func NewClog() *Clog {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(handler)
	return &Clog{
		logger: handler,
	}
}

// Info logs an info message
func (c *Clog) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Info(msg, keysAndValues...)
}

// Error logs an error message
func (c *Clog) Error(msg string, keysAndValues ...interface{}) {
	c.logger.Error(msg, keysAndValues...)
}

// Debug logs a debug message
func (c *Clog) Debug(msg string, keysAndValues ...interface{}) {
	c.logger.Debug(msg, keysAndValues...)
}

// Warn logs a warning message
func (c *Clog) Warn(msg string, keysAndValues ...interface{}) {
	c.logger.Warn(msg, keysAndValues...)
}
