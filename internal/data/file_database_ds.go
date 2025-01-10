package data

import (
	"context"
	"time"

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

func (d *fileDatabaseDs) UpdateFile(ctx context.Context, file *domain.File) (*domain.File, error) {
	return nil, nil
}

func (d *fileDatabaseDs) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return nil
}
