package handler

import (
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
)

// Handler for the me endpoint
func Me(srv service.AuthenticationService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			res, err := srv.Me(r.Context())
			if err != nil {
				switch err.Error() {
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			response.OK(w, r, res)
		},
	)
}
