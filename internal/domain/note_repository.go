package domain

import (
	"context"

	"github.com/google/uuid"
)

type NoteRepository interface {
	ListNote(ctx context.Context) (*[]Note, error)
	GetNote(ctx context.Context, id uuid.UUID) (*Note, error)
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	UpdateNote(ctx context.Context, note *Note) (*Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error
}
