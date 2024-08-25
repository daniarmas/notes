package server

import (
	"encoding/json"
	"net/http"

	"github.com/daniarmas/notes/internal/config"
)

// Instantiate a net/http server...
func NewServer(
	config *config.Configuration,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
	)
	var handler http.Handler = mux
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	config *config.Configuration,
) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /health", handleHealthCheck)
}

// healthcheck handler
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// HealthCheckResponse represents the structure of the health check response
	type HealthCheckResponse struct {
		Status string `json:"status"`
	}
	response := HealthCheckResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
