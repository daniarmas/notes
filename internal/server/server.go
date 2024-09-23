package server

import (
	"net/http"

	"github.com/daniarmas/notes/internal/server/handler"
	"github.com/daniarmas/notes/internal/server/middleware"
)

type HandleFunc struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

type Server struct {
	Mux        *http.ServeMux
	HttpServer *http.Server
}

// NewServer creates and configures a new HTTP server with the specified address.
func NewServer(addr string) *Server {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Add routes

	// Not found
	mux.HandleFunc("/", handler.NotFoundHandler)
	// Health check
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	// OpenAPI specification
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})
	// Swagger UI
	mux.Handle("GET /doc/", http.StripPrefix("/doc", http.FileServer(http.Dir("docs/swaggerui/dist"))))

	var handler http.Handler = mux
	// Add logging middleware
	handler = middleware.LoggingMiddleware(handler)
	// Create the HTTP server
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &Server{
		Mux:        mux,
		HttpServer: httpServer,
	}
}

// AddRoutes adds routes to the server
func (s *Server) AddRoutes(
	handlers []HandleFunc,
) {
	for _, h := range handlers {
		s.Mux.HandleFunc(h.Pattern, h.Handler)
	}
}
