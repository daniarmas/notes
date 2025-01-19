package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type NoteDatabaseDs interface {
	ListNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	ListTrashNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]Note, error)
	GetNote(ctx context.Context, id uuid.UUID) (*Note, error)
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	UpdateNote(ctx context.Context, note *Note) (*Note, error)
	RestoreNote(ctx context.Context, id uuid.UUID) (*Note, error)
	HardDeleteNote(ctx context.Context, id uuid.UUID) error
	SoftDeleteNote(ctx context.Context, id uuid.UUID) error
}
