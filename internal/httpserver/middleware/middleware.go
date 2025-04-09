package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/daniarmas/clogg"
	cmiddleware "github.com/daniarmas/http/middleware"
	"github.com/daniarmas/http/response"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/utils"
	"github.com/google/uuid"
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
func SetUserInContext(jwtDatasource domain.JwtDatasource) cmiddleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header from the request
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Split the header to get the token
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					response.Unauthorized(w, r, "Authorization header format is invalid. It must be in the format: 'Bearer {token}'.", nil)
					return
				}
				token := parts[1]
				jwtMetadata := domain.JWTMetadata{Token: token}
				err := jwtDatasource.ParseJWT(&jwtMetadata)
				if err != nil {
					switch err.Error() {
					case "Token is expired":
						response.Unauthorized(w, r, "Authorization token has expired. Please log in again to continue.", nil)
						return
					case "signature is invalid":
						response.Unauthorized(w, r, "Authorization token signature is invalid. Please provide a valid token.", nil)
						return
					case "token contains an invalid number of segments":
						response.Unauthorized(w, r, "Authorization token provided is invalid. Please provide a valid token.", nil)
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
		clogg.Info(
			ctx,
			"HTTP request",
			clogg.String("method", r.Method),
			clogg.String("path", r.URL.Path),
			clogg.Int("status", rw.statusCode),
			clogg.String("user_agent", r.Header.Get("User-Agent")),
			clogg.String("request_id", requestID),
		)
	})
}

func LoggedOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := domain.GetUserIdFromContext(r.Context())
		if userId == uuid.Nil {
			response.Unauthorized(w, r, "User is not logged in. Please log in to access this resource.", nil)
			return
		}
		// Call the next handler
		h.ServeHTTP(w, r)
	})
}
