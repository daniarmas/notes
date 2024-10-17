package httpserver

import (
	"github.com/daniarmas/notes/internal/httpserver/handler"
	"github.com/daniarmas/notes/internal/service"
)

func Routes(authenticationService service.AuthenticationService) []HandleFunc {
	return []HandleFunc{
		// Default routes

		// Not found
		{Pattern: "/", Handler: handler.NotFoundHandler},
		// Health check
		{Pattern: "GET /health", Handler: handler.HealthCheckHandler},
		// OpenAPI specification
		{Pattern: "GET /swagger.json", Handler: handler.OpenApiHanlder},

		// Authentication
		{Pattern: "POST /sign-in", Handler: handler.SignIn(authenticationService)},
	}
}
