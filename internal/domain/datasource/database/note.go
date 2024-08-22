package database

import (
	"context"

	"github.com/daniarmas/notes/internal/entity"
	"github.com/google/uuid"
)

type NoteDatabaseDs interface {
	ListNote(ctx context.Context) (*[]entity.Note, error)
	GetNote(ctx context.Context, id uuid.UUID) (*entity.Note, error)
	CreateNote(ctx context.Context, user *entity.Note) (*entity.Note, error)
	UpdateNote(ctx context.Context, user *entity.Note) (*entity.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error
}
