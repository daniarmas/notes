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
		{Pattern: "GET /openapi.yaml", Handler: handler.OpenApiHanlder},
		// Swagger UI
		{Pattern: "GET /doc/", Handler: handler.SwaggerUiHandler},

		// Authentication
		{Pattern: "POST /sign-in", Handler: handler.SignIn(authenticationService)},
	}
}
