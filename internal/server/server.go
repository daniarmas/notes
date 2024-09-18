package server

import (
	"net/http"

	"github.com/daniarmas/notes/internal/server/controller"
	"github.com/daniarmas/notes/internal/server/middleware"
	"github.com/daniarmas/notes/internal/service"
)

// NewServer creates a new HTTP server
func NewServer(
	authenticationService service.AuthenticationService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		authenticationService,
	)
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	return handler
}

// addRoutes adds the routes to the HTTP server
func addRoutes(
	mux *http.ServeMux,
	authenticationService service.AuthenticationService,
) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /health", controller.HandleHealthCheck)
	mux.HandleFunc("POST /sign-in", controller.HandleSignIn(authenticationService))
}
