package service

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type CreateNoteResponse struct {
	Note *domain.Note `json:"note"`
}

type NoteService interface {
	CreateNote(ctx context.Context, title string, content string) (*CreateNoteResponse, error)
	ListNotesByUser(ctx context.Context, cursor time.Time) (*[]domain.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error
	UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error)
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

func (s *noteService) UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	note, err := s.NoteRepository.UpdateNote(ctx, note)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("note not found")
		}
	}
	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID) error {
	err := s.NoteRepository.DeleteNote(ctx, id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return errors.New("note not found")
		}
	}
	return nil
}
