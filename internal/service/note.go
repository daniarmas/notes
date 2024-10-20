package service

import (
	"context"
	"time"

	"github.com/daniarmas/notes/internal/domain"
)

type CreateNoteResponse struct {
	Note *domain.Note `json:"note"`
}

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string) (*CreateNoteResponse, error)
	ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error)
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
		UserId:  domain.GetUserIdFromContext(ctx),
		Title:   title,
		Content: content,
	}
	note, err := s.NoteRepository.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}
	return &CreateNoteResponse{Note: note}, nil
}

func (s *noteService) ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error) {
	// Get the user ID from the context
	userId := domain.GetUserIdFromContext(ctx)

	notes, err := s.NoteRepository.ListNotesByUser(ctx, userId, cursor)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
