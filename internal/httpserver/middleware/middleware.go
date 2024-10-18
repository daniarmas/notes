package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver/response"
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
func SetUserInContext(next http.Handler, jwtDatasource domain.JwtDatasource) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header from the request
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Split the header to get the token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, r, "authorization header format must be Bearer {token}", nil)
				return
			}
			token := parts[1]
			jwtMetadata := domain.JWTMetadata{Token: token}
			err := jwtDatasource.ParseJWT(&jwtMetadata)
			if err != nil {
				switch err.Error() {
				case "Token is expired":
					response.Unauthorized(w, r, "authorization token expired", nil)
					return
				case "signature is invalid":
					response.Unauthorized(w, r, "authorization token signature is invalid", nil)
					return
				case "token contains an invalid number of segments":
					response.Unauthorized(w, r, "authorization token is invalid", nil)
					return
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			// Set the user in the context
			ctx := domain.SetUserInContext(r.Context(), jwtMetadata.UserId)

			// Call the next handler with the modified context
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Call the next handler
			next.ServeHTTP(w, r)
		}
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
