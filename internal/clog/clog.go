package clog

import (
	"context"
	"log/slog"
	"os"

	"github.com/daniarmas/notes/internal/server/utils"
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
func (c *Clog) Info(ctx context.Context, msg string, err error) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg,
		slog.String("error", err.Error()),
		slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
	)
}

// Error logs an error message
func (c *Clog) Error(ctx context.Context, msg string, err error) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelError,
		msg,
		slog.String("error", err.Error()),
		slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
	)
}

// Debug logs a debug message
func (c *Clog) Debug(ctx context.Context, msg string, err error) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		msg,
		slog.String("error", err.Error()),
		slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
	)
}

// Warn logs a warning message
func (c *Clog) Warn(ctx context.Context, msg string, err error) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelWarn,
		msg,
		slog.String("error", err.Error()),
		slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
	)
}
