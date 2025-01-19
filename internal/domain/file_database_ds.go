package domain

import (
	"context"

	"github.com/google/uuid"
)

type FileDatabaseDs interface {
	ListFilesByNotesIds(ctx context.Context, noteId []uuid.UUID) (*[]File, error)
	ListFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]File, error)
	CreateFile(ctx context.Context, file *File) (*File, error)
	UpdateFileByOriginalId(ctx context.Context, originalFileId, processFileId string) (*File, error)
	HardDeleteFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]File, error)
}
