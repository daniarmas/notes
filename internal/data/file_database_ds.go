package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type fileDatabaseDs struct {
	queries *database.Queries
}

func NewFileDatabaseDs(queries *database.Queries) domain.FileDatabaseDs {
	return &fileDatabaseDs{
		queries: queries,
	}
}

func (d *fileDatabaseDs) ListFilesByNote(ctx context.Context, note_id uuid.UUID) (*[]domain.File, error) {
	return nil, nil
}

func (d *fileDatabaseDs) CreateFile(ctx context.Context, file *domain.File) (*domain.File, error) {
	// Get current time
	timeNow := time.Now().UTC()

	res, err := d.queries.CreateFile(ctx, database.CreateFileParams{
		NoteID:       file.NoteId,
		OriginalFile: file.OriginalFile,
		CreateTime:   timeNow,
		UpdateTime:   timeNow,
	})
	if err != nil {
		return nil, err
	}
	return &domain.File{
		Id:           res.ID,
		NoteId:       res.NoteID,
		OriginalFile: res.OriginalFile,
		CreateTime:   res.CreateTime,
		UpdateTime:   res.UpdateTime,
		DeleteTime:   res.DeleteTime.Time,
	}, nil
}

func (d *fileDatabaseDs) UpdateFileByOriginalId(ctx context.Context, originalFileId, processFileId string) (*domain.File, error) {
	// Get current time
	timeNow := time.Now().UTC()

	res, err := d.queries.UpdateFileByOriginalId(ctx, database.UpdateFileByOriginalIdParams{OriginalFile: originalFileId, ProcessedFile: sql.NullString{String: processFileId, Valid: true}, UpdateTime: timeNow})
	if err != nil {
		clog.Error(ctx, "error updating file by original id", err)
		return nil, err
	}

	file := parseToDomain(res)

	return file, nil
}

func (d *fileDatabaseDs) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return nil
}

// ParseToDomain parses a file from the database to a domain.File
func parseToDomain(f database.File) *domain.File {
	// Parse UUIDs and handle potential errors
	id, _ := uuid.Parse(f.ID.String())
	noteId, _ := uuid.Parse(f.NoteID.String())

	// Return the parsed domain.File
	return &domain.File{
		Id:            id,
		NoteId:        noteId,
		OriginalFile:  f.OriginalFile,
		ProcessedFile: f.ProcessedFile.String,
		CreateTime:    f.CreateTime,
		UpdateTime:    f.UpdateTime,
		DeleteTime:    f.DeleteTime.Time,
	}
}
