package controller

import (
	"encoding/json"
	"net/http"

	"github.com/daniarmas/notes/internal/server/response"
	"github.com/daniarmas/notes/internal/server/validate"
	"github.com/daniarmas/notes/internal/service"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignInRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.Email == "" {
		errors["email"] = "email is required"
	} else {
		validate.ValidateEmail(&errors, r.Email)
	}
	if r.Password == "" {
		errors["password"] = "password is required"
	}
	return errors
}

func HandleSignIn(srv service.AuthenticationService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Parse the request body into a SignInRequest struct
			var req SignInRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				response.BadRequestHandler(w, r, "Invalid JSON request", nil)
				return
			}
			defer r.Body.Close()

			// Validate the request and return an InvalidRequestDataError if there are any errors
			if errors := req.Validate(); len(errors) > 0 {
				response.BadRequestHandler(w, r, "Invalid request data", nil)
				return
			}

			res, err := srv.SignIn(r.Context(), req.Email, req.Password)
			if err != nil {
				switch err.Error() {
				case "invalid credentials":
					response.UnauthorizedHandler(w, r, "Invalid credentials", nil)
					return
				default:
					response.InternalServerErrorHandler(w, r)
					return
				}
			}
			response.StatusOk(w, r, res)
		},
	)
}
