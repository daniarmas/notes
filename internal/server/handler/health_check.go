package handler

import (
	"net/http"
)

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler is the handler for the health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	res := HealthCheckResponse{Status: "healthy"}
	StatusOk(w, r, res)
}
