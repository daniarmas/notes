package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rs/xid"
)

// responseWriter is a custom http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type contextKey string

const RequestIDKey contextKey = "request-id"

// LoggingMiddleware injects the request ID into the context and logs the request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID
		requestID := xid.New().String()

		// Add the request ID to the context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Create a custom response writer to capture the status code
		rw := &responseWriter{ResponseWriter: w}

		// Call the next handler with the modified context
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Log the request details, optionally retrieving the request ID from the context
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.String("user_agent", r.Header.Get("User-Agent")),
			slog.String("request_id", requestID), // You can also retrieve it from ctx if necessary
		)
	})
}
