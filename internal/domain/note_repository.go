package domain

import (
	"context"
	"github.com/daniarmas/notes/internal/clog"
)

type NoteRepository interface {
	// ListNote(ctx context.Context) (*[]Note, error)
	// GetNote(ctx context.Context, id uuid.UUID) (*Note, error)
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	// UpdateNote(ctx context.Context, note *Note) (*Note, error)
	// DeleteNote(ctx context.Context, id uuid.UUID) error
}

type noteRepository struct {
	NoteCacheDs    NoteCacheDs
	NoteDatabaseDs NoteDatabaseDs
}

func NewNoteRepository(noteCacheDs *NoteCacheDs, noteDatabaseDs *NoteDatabaseDs) NoteRepository {
	return &noteRepository{
		NoteCacheDs:    *noteCacheDs,
		NoteDatabaseDs: *noteDatabaseDs,
	}
}

func (n *noteRepository) CreateNote(ctx context.Context, note *Note) (*Note, error) {
	// Save the note on the database
	note, err := n.NoteDatabaseDs.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}
	// Cache the note
	err = n.NoteCacheDs.CreateNote(ctx, note)
	if err != nil {
		clog.Info(
			ctx,
			err.Error(),
			nil,
		)
	}
	return note, nil
}
