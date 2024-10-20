package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver/response"
	"github.com/daniarmas/notes/internal/service"
	"github.com/daniarmas/notes/internal/utils"
	"github.com/google/uuid"
)

// Represents the structure of the create note request
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Represent the structure of the list notes response
type ListNotesResponse struct {
	Notes  *[]domain.Note `json:"notes"`
	Cursor time.Time      `json:"cursor"`
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

			response.OK(w, r, res)
		},
	)
}

// Handler for the list notes endpoint
func ListNotesByUser(srv service.NoteService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Get the cursor from the query parameters
			cursorQueryParam := r.URL.Query().Get("cursor")
			// parse the cursor query parameter
			cursor, err := utils.ParseTime(cursorQueryParam)
			if err != nil && cursorQueryParam != "" {
				msg := "Invalid time format for the cursor query parameter. Must use RFC3339 format"
				response.BadRequest(w, r, &msg, nil)
				return
			}

			// If the cursor is zero, set it to the current time
			if cursor.IsZero() {
				cursor = time.Now().UTC()
			}

			notes, err := srv.ListNotesByUser(r.Context(), cursor)
			if err != nil {
				switch err.Error() {
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			// Get the next cursor
			notesSlice := *notes
			var nextCursor time.Time
			if len(notesSlice) > 0 {
				nextCursor = notesSlice[len(notesSlice)-1].CreateTime
			} else {
				// Handle the case where notesSlice is empty
				nextCursor = time.Now().UTC()
			}

			res := ListNotesResponse{
				Notes:  notes,
				Cursor: nextCursor,
			}
			response.OK(w, r, res)
		},
	)
}

// Handler for the delete note endpoint
func DeleteNote(srv service.NoteService) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Get the note ID from the URL path
			idPathParam := r.PathValue("id")
			id, err := uuid.Parse(idPathParam)
			if err != nil {
				msg := "Provided ID path parameter is invalid. It must be a valid UUID."
				response.BadRequest(w, r, &msg, nil)
				return
			}

			err = srv.DeleteNote(r.Context(), id)
			if err != nil {
				switch err.Error() {
				case "note not found":
					response.NotFound(w, r, "")
					return
				default:
					response.InternalServerError(w, r)
					return
				}
			}

			response.NotContent(w, r)
		},
	)
}
