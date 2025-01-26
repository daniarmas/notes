package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/google/uuid"
)

type NoteRepository interface {
	ListNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	ListTrashNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	// GetNote(ctx context.Context, id uuid.UUID) (*Note, error)
	CreateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error)
	RestoreNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*Note, error)
	UpdateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error)
	DeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID, isHard bool) error
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

func (n *noteRepository) CreateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error) {
	// Save the note on the database
	note, err := n.NoteDatabaseDs.CreateNote(ctx, tx, note)
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

func (n *noteRepository) ListNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error) {
	// Fetch the notes from the database
	notes, err := n.NoteDatabaseDs.ListNotesByUser(ctx, user_id, cursor)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (n *noteRepository) ListTrashNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error) {
	// Fetch the notes from the database
	notes, err := n.NoteDatabaseDs.ListTrashNotesByUser(ctx, user_id, cursor)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (n *noteRepository) UpdateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error) {
	// Update the note on the database
	note, err := n.NoteDatabaseDs.UpdateNote(ctx, tx, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (n *noteRepository) RestoreNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*Note, error) {
	// Update the note on the database
	note, err := n.NoteDatabaseDs.RestoreNote(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (n *noteRepository) DeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID, isHard bool) error {
	if isHard {
		// Hard delete the note from the database
		err := n.NoteDatabaseDs.HardDeleteNote(ctx, tx, id)
		if err != nil {
			return err
		}
	} else {
		// Soft delete the note from the database
		err := n.NoteDatabaseDs.SoftDeleteNote(ctx, tx, id)
		if err != nil {
			return err
		}
	}
	return nil
}
