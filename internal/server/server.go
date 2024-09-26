package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/daniarmas/notes/internal/server/middleware"
	"github.com/daniarmas/notes/internal/server/response"
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
func NewServer(addr string, handlers []HandleFunc) *Server {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Add routes
	// Business routes
	for _, h := range handlers {
		mux.HandleFunc(h.Pattern, h.Handler)
	}
	// Not found
	mux.HandleFunc("/", NotFoundHandler)
	// Health check
	mux.HandleFunc("GET /health", HealthCheckHandler)
	// OpenAPI specification
	mux.HandleFunc("/openapi.yaml", OpenApiHanlder)
	// Swagger UI
	mux.Handle("GET /doc/", http.StripPrefix("/doc", http.FileServer(http.Dir("third_party/swaggerui/dist"))))

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

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler handles HTTP requests for checking the health status of the server.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	res := HealthCheckResponse{Status: "healthy"}
	response.StatusOk(w, r, res)
}

// NotFoundHandler handles requests to non-existent resources.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Determine the scheme
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	docUrl := fmt.Sprintf("%s://%s/doc", scheme, r.Host)

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "" && (strings.Contains(acceptHeader, "text/html") || strings.Contains(acceptHeader, "application/xhtml+xml")) {
		http.Redirect(w, r, docUrl, http.StatusFound)
		return
	}

	msg := fmt.Sprintf("Resource not found. Please refer to the documentation for further details. Doc: %s", docUrl)
	response.NotFound(w, r, msg)
}

// OpenApiHanlder handles requests for the OpenAPI specification.
func OpenApiHanlder(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "api/openapi.yaml")
}
