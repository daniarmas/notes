package service

import (
	"context"

	"github.com/daniarmas/notes/internal/domain"
)

type CreateNoteResponse struct {
	Note domain.Note `json:"note"`
}

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string) (*CreateNoteResponse, error)
}


