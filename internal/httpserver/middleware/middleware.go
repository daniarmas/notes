package middleware

import (
	"context"
	"net/http"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/domain"
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

// SetUserInContext is a middleware that sets the user in the context
func SetUserInContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the request context
		user := domain.GetUserFromContext(r.Context())

		// Set the user in the context
		ctx := domain.SetUserInContext(r.Context(), user)

		// Call the next handler with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AllowCORS is a middleware that sets the CORS headers
func AllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request is an OPTIONS request, return immediately
		if r.Method == "OPTIONS" {
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
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
