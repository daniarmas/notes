package service

import (
	"context"

	"github.com/daniarmas/notes/internal/domain"
)

type CreateNoteResponse struct {
	Note *domain.Note `json:"note"`
}

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string) (*CreateNoteResponse, error)
}

type noteService struct {
	NoteRepository domain.NoteRepository
}

func NewNoteService(noteRepository domain.NoteRepository) NoteService {
	return &noteService{
		NoteRepository: noteRepository,
	}
}

func (s *noteService) CreateNote(ctx context.Context, title string, content string) (*CreateNoteResponse, error) {
	note := &domain.Note{
		Title:   title,
		Content: content,
	}
	note, err := s.NoteRepository.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}
	return &CreateNoteResponse{Note: note}, nil
}
