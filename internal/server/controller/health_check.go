package controller

import (
	"net/http"

	"github.com/daniarmas/notes/internal/server/response"
)

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HandleHealthCheck is the handler for the health check endpoint
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := HealthCheckResponse{Status: "healthy"}
	response.StatusOk(w, r, res)
}
