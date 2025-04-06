package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

// Parses a file from the database to a domain.File
func parseFromDatabaseToDomain(f database.File) domain.File {
	return domain.File{
		Id:            f.ID,
		NoteId:        f.NoteID,
		OriginalFile:  f.OriginalFile,
		ProcessedFile: f.ProcessedFile.String,
		CreateTime:    f.CreateTime,
		UpdateTime:    f.UpdateTime,
		DeleteTime:    f.DeleteTime.Time,
	}
}

type fileDatabaseDs struct {
	queries *database.Queries
}

func NewFileDatabaseDs(queries *database.Queries) domain.FileDatabaseDs {
	return &fileDatabaseDs{
		queries: queries,
	}
}

func (d *fileDatabaseDs) ListFilesByNotesIds(ctx context.Context, noteId []uuid.UUID) (*[]domain.File, error) {
	res, err := d.queries.ListFilesByNotesIds(ctx, noteId)
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.File, 0, len(res))
	for _, file := range res {
		response = append(response, parseFromDatabaseToDomain(file))
	}
	return &response, nil
}

func (d *fileDatabaseDs) ListFilesByNoteId(ctx context.Context, noteId uuid.UUID) (*[]domain.File, error) {
	res, err := d.queries.ListFileByNoteId(ctx, noteId)
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.File, 0, len(res))
	for _, file := range res {
		response = append(response, parseFromDatabaseToDomain(file))
	}
	return &response, nil

}

func (d *fileDatabaseDs) CreateFile(ctx context.Context, tx *sql.Tx, file *domain.File) (*domain.File, error) {
	// Get current time
	timeNow := time.Now().UTC()

	res, err := d.queries.WithTx(tx).CreateFile(ctx, database.CreateFileParams{
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

func (d *fileDatabaseDs) UpdateFileByOriginalId(ctx context.Context, tx *sql.Tx, originalFileId, processFileId string) (*domain.File, error) {
	// Get current time
	timeNow := time.Now().UTC()

	res, err := d.queries.WithTx(tx).UpdateFileByOriginalId(ctx, database.UpdateFileByOriginalIdParams{OriginalFile: originalFileId, ProcessedFile: sql.NullString{String: processFileId, Valid: true}, UpdateTime: timeNow})
	if err != nil {
		clogg.Error(ctx, "error updating file by original id")
		return nil, err
	}

	file := parseToDomain(res)

	return file, nil
}

func (d *fileDatabaseDs) HardDeleteFilesByNoteId(ctx context.Context, tx *sql.Tx, noteId uuid.UUID) (*[]domain.File, error) {
	res, err := d.queries.WithTx(tx).HardDeleteFilesByNoteId(ctx, noteId)
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.File, 0, len(res))
	for _, file := range res {
		response = append(response, parseFromDatabaseToDomain(file))
	}
	return &response, nil
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
