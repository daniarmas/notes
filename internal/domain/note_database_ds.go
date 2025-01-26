package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type NoteDatabaseDs interface {
	ListNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	ListTrashNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	GetNote(ctx context.Context, id uuid.UUID) (*Note, error)
	CreateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error)
	UpdateNote(ctx context.Context, tx *sql.Tx, note *Note) (*Note, error)
	RestoreNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*Note, error)
	HardDeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
	SoftDeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
}
