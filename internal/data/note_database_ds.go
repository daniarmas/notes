package data

import (
	"context"
	"database/sql"

	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type noteDatabaseDs struct {
	queries *database.Queries
}

func NewNoteDatabaseDs(queries *database.Queries) domain.NoteDatabaseDs {
	return &noteDatabaseDs{
		queries: queries,
	}
}

func (d *noteDatabaseDs) CreateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	res, err := d.queries.CreateNote(ctx, database.CreateNoteParams{
		UserID:          note.UserId,
		Title:           sql.NullString{String: note.Title},
		Content:         sql.NullString{String: note.Content},
		BackgroundColor: sql.NullString{String: note.BackgroundColor},
	})
	if err != nil {
		return nil, err
	}
	return &domain.Note{
		Id:              res.ID,
		UserId:          res.UserID,
		Title:           res.Title.String,
		Content:         res.Content.String,
		BackgroundColor: res.BackgroundColor.String,
		CreateTime:      res.CreateTime,
	}, nil
}

func (d *noteDatabaseDs) ListNotesByUserId(ctx context.Context, user_id uuid.UUID) (*[]domain.Note, error) {
	res, err := d.queries.ListNotesByUserId(ctx, user_id)
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.Note, 0, len(res))
	for _, note := range res {
		response = append(response, domain.Note{
			Id:              note.ID,
			UserId:          note.UserID,
			Title:           note.Title.String,
			Content:         note.Content.String,
			BackgroundColor: note.BackgroundColor.String,
			CreateTime:      note.CreateTime,
			UpdateTime:      note.UpdateTime.Time,
			DeleteTime:      note.DeleteTime.Time,
		})
	}
	return &response, nil
}

func (d *noteDatabaseDs) GetNote(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	return nil, nil
}

func (d *noteDatabaseDs) UpdateNote(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	return nil, nil
}
func (d *noteDatabaseDs) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return nil
}