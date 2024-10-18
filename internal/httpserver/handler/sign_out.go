package handler

import (
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
)

// Handler for the sign-out endpoint
func SignOut(srv service.AuthenticationService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			err := srv.SignOut(r.Context())
			if err != nil {
				switch err.Error() {
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			response.StatusOk(w, r, nil)
		},
	)
}