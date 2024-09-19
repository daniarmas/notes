package clog

import (
	"context"
	"log/slog"
	"os"

	"github.com/daniarmas/notes/internal/server/utils"
)

// Clog is a custom logger that wraps the slog package
type Clog struct {
	Logger *slog.Logger
}

// NewClog creates a new Clog instance
func NewClog() *Clog {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(handler)
	return &Clog{
		Logger: handler,
	}
}

// Info logs an info message
func Info(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	attributes := append(attrs, slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)))
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg,
		attributes...,
	)
}

// Error logs an error message
func Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	attributes := append(attrs, slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)))
	slog.LogAttrs(
		context.Background(),
		slog.LevelError,
		msg,
		attributes...,
	)
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	attributes := append(attrs, slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)))
	slog.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		msg,
		attributes...,
	)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	attributes := append(attrs, slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)))
	slog.LogAttrs(
		context.Background(),
		slog.LevelWarn,
		msg,
		attributes...,
	)
}

func String(key string, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

func Uint64(key string, value uint64) slog.Attr {
	return slog.Uint64(key, value)
}

func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}
