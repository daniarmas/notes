package handler

import (
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
	"github.com/daniarmas/notes/internal/validate"
)

// Represents the structure of the sign-in request
type SignOutRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validates the sign-in request
func (r SignOutRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.Email == "" {
		errors["email"] = "field required"
	} else {
		validate.ValidateEmail(&errors, r.Email)
	}
	if r.Password == "" {
		errors["password"] = "field required"
	}
	return errors
}

// Handler for the sign-in endpoint
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
