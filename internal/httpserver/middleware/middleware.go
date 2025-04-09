package middleware

import (
	"net/http"
	"strings"

	cmiddleware "github.com/daniarmas/http/middleware"
	"github.com/daniarmas/http/response"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
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
