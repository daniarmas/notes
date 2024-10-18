package handler

import (
	"encoding/json"
	"net/http"

	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
)

// Represents the structure of the create note request
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Validates the create note request
func (r CreateNoteRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.Title == "" {
		errors["title"] = "field required"
	}
	return errors
}

// Handler for the sign-in endpoint
func CreateNote(srv service.NoteService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Parse the request body into a CreateNoteRequest struct
			var req CreateNoteRequest
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

			res, err := srv.CreateNote(r.Context(), req.Title, req.Content)
			if err != nil {
				switch err.Error() {
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			response.StatusOk(w, r, res)
		},
	)
}
