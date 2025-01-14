package domain

import (
	"context"

	"github.com/google/uuid"
)

type FileDatabaseDs interface {
	ListFilesByNote(ctx context.Context, note_id uuid.UUID) (*[]File, error)
	CreateFile(ctx context.Context, file *File) (*File, error)
	UpdateFileByOriginalId(ctx context.Context, originalFileId, processFileId string) (*File, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error
}
