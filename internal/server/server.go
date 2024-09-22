package server

import (
	"net"
	"net/http"

	"github.com/daniarmas/notes/internal/server/handler"
	"github.com/daniarmas/notes/internal/server/middleware"
	"github.com/daniarmas/notes/internal/service"
)

// NewServer creates a new HTTP server
func NewServer(
	authenticationService service.AuthenticationService,
) *http.Server {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		authenticationService,
	)
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	return &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", "8080"),
		Handler: handler,
	}
}

// addRoutes adds the routes to the HTTP server
func addRoutes(
	mux *http.ServeMux,
	authenticationService service.AuthenticationService,
) {
	mux.HandleFunc("/", handler.NotFoundHandler)
	// Health check
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	// OpenAPI specification
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})
	// Swagger UI
	mux.Handle("GET /doc/", http.StripPrefix("/doc", http.FileServer(http.Dir("docs/swaggerui/dist"))))

	// Notes

	// Authentication
	mux.HandleFunc("POST /sign-in", handler.SignInHandler(authenticationService))
}
