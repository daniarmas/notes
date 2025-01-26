package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type FileDatabaseDs interface {
	ListFilesByNotesIds(ctx context.Context, noteId []uuid.UUID) (*[]File, error)
	ListFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]File, error)
	CreateFile(ctx context.Context, tx *sql.Tx, file *File) (*File, error)
	UpdateFileByOriginalId(ctx context.Context, tx *sql.Tx, originalFileId, processFileId string) (*File, error)
	HardDeleteFilesByNoteId(ctx context.Context, tx *sql.Tx, noteId uuid.UUID) (*[]File, error)
}
