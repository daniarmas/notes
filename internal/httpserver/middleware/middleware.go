package middleware

import (
	"context"
	"net/http"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/utils"
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

// LoggingMiddleware injects the request ID into the context and logs the request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID
		requestID := xid.New().String()

		// Add the request ID to the context
		ctx := context.WithValue(r.Context(), utils.RequestIDKey, requestID)

		// Create a custom response writer to capture the status code
		rw := &responseWriter{ResponseWriter: w}

		// Call the next handler with the modified context
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Log the request details, optionally retrieving the request ID from the context
		// clog.Debug(ctx, "HTTP request", nil)
		clog.Info(
			ctx,
			"HTTP request",
			nil,
			clog.String("method", r.Method),
			clog.String("path", r.URL.Path),
			clog.Int("status", rw.statusCode),
			clog.String("user_agent", r.Header.Get("User-Agent")),
			clog.String("request_id", requestID),
		)
	})
}