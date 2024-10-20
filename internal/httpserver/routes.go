package httpserver

import (
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/handler"
	"github.com/daniarmas/notes/internal/httpserver/middleware"
	"github.com/daniarmas/notes/internal/service"
)

func Routes(authenticationService service.AuthenticationService, noteService service.NoteService) []HandleFunc {
	return []HandleFunc{
		// Default routes

		// Not found
		{Pattern: "/", Handler: handler.NotFoundHandler},
		// Health check
		{Pattern: "GET /health", Handler: handler.HealthCheckHandler},
		// OpenAPI specification
		{Pattern: "GET /swagger.json", Handler: handler.OpenApiHanlder},

		// Authentication
		{Pattern: "GET /me", Handler: middleware.LoggedOnly(handler.Me(authenticationService)).(http.HandlerFunc)},
		{Pattern: "POST /sign-in", Handler: handler.SignIn(authenticationService)},
		{Pattern: "POST /sign-out", Handler: middleware.LoggedOnly(handler.SignOut(authenticationService)).(http.HandlerFunc)},

		// Note
		{Pattern: "GET /note", Handler: middleware.LoggedOnly(handler.ListNotesByUser(noteService)).(http.HandlerFunc)},
		{Pattern: "POST /note", Handler: middleware.LoggedOnly(handler.CreateNote(noteService)).(http.HandlerFunc)},
	}
}
