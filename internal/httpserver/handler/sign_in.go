package handler

import (
	"encoding/json"
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
	"github.com/daniarmas/notes/internal/validate"
)

// Represents the structure of the sign-in request
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validates the sign-in request
func (r SignInRequest) Validate() map[string]string {
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
func SignIn(srv service.AuthenticationService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Parse the request body into a SignInRequest struct
			var req SignInRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				msg := "Invalid JSON request"
				response.BadRequest(w, r, &msg, nil)
				return
			}
			defer r.Body.Close()

			// Validate the request and return an InvalidRequestDataError if there are any errors
			if errors := req.Validate(); len(errors) > 0 {
				response.BadRequest(w, r, nil, errors)
				return
			}

			res, err := srv.SignIn(r.Context(), req.Email, req.Password)
			if err != nil {
				switch err.Error() {
				case "invalid credentials":
					response.Unauthorized(w, r, "Invalid credentials", nil)
					return
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			response.OK(w, r, res)
		},
	)
}
