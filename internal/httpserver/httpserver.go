package httpserver

import (
	"net"
	"net/http"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver/middleware"
	"github.com/daniarmas/notes/internal/service"
)

type HandleFunc struct {
	Pattern string
	Handler http.HandlerFunc
}

type Server struct {
	Mux        *http.ServeMux
	HttpServer *http.Server
}

// NewServer creates and configures a new HTTP server with the specified address.
func NewServer(authenticationService service.AuthenticationService, noteService service.NoteService, jwtDatasource domain.JwtDatasource) *Server {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Routes
	routes := Routes(authenticationService, noteService)
	for _, h := range routes {
		mux.HandleFunc(h.Pattern, h.Handler)
	}

	var handler http.Handler = mux
	// Add middlewares
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.AllowCORS(handler)
	handler = middleware.SetUserInContext(handler, jwtDatasource)

	// Create the HTTP server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", "8080"),
		Handler: handler,
	}
	return &Server{
		Mux:        mux,
		HttpServer: httpServer,
	}
}
