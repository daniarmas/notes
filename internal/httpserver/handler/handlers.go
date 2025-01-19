package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/daniarmas/notes/internal/httpserver/response"
)

// OpenApiHanlder handles requests for the OpenAPI specification.
func OpenApiHanlder(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "api/openapi-spec/swagger.json")
}

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler handles HTTP requests for checking the health status of the server.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	res := HealthCheckResponse{Status: "healthy"}
	response.OK(w, r, res)
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
