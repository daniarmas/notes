package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
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

func (d *noteDatabaseDs) CreateNote(ctx context.Context, tx *sql.Tx, note *domain.Note) (*domain.Note, error) {
	// Get current time
	timeNow := time.Now().UTC()

	res, err := d.queries.WithTx(tx).CreateNote(ctx, database.CreateNoteParams{
		UserID:     note.UserId,
		Title:      sql.NullString{String: note.Title, Valid: true},
		Content:    sql.NullString{String: note.Content, Valid: true},
		CreateTime: timeNow,
		UpdateTime: timeNow,
	})
	if err != nil {
		return nil, err
	}
	return &domain.Note{
		Id:         res.ID,
		UserId:     res.UserID,
		Title:      res.Title.String,
		Content:    res.Content.String,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		DeleteTime: res.DeleteTime.Time,
	}, nil
}

func (d *noteDatabaseDs) ListNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]domain.Note, error) {
	res, err := d.queries.ListNotesByUserId(ctx, database.ListNotesByUserIdParams{UserID: user_id, UpdateTime: cursor})
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.Note, 0, len(res))
	for _, note := range res {
		response = append(response, domain.Note{
			Id:         note.ID,
			UserId:     note.UserID,
			Title:      note.Title.String,
			Content:    note.Content.String,
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
			DeleteTime: note.DeleteTime.Time,
		})
	}
	return &response, nil
}

func (d *noteDatabaseDs) ListTrashNotesByUser(ctx context.Context, user_id uuid.UUID, cursor time.Time) (*[]domain.Note, error) {
	res, err := d.queries.ListTrashNotesByUserId(ctx, database.ListTrashNotesByUserIdParams{UserID: user_id, DeleteTime: sql.NullTime{Time: cursor, Valid: true}})
	if err != nil {
		return nil, err
	}
	// Preallocate slice with the length of the result set
	response := make([]domain.Note, 0, len(res))
	for _, note := range res {
		response = append(response, domain.Note{
			Id:         note.ID,
			UserId:     note.UserID,
			Title:      note.Title.String,
			Content:    note.Content.String,
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
			DeleteTime: note.DeleteTime.Time,
		})
	}
	return &response, nil
}

func (d *noteDatabaseDs) GetNote(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	return nil, nil
}

func (d *noteDatabaseDs) UpdateNote(ctx context.Context, tx *sql.Tx, note *domain.Note) (*domain.Note, error) {
	res, err := d.queries.WithTx(tx).UpdateNoteById(ctx, database.UpdateNoteByIdParams{
		ID:         note.Id,
		Column2:    note.Title,
		Column3:    note.Content,
		UpdateTime: time.Now().UTC(),
	})
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, err
		}
	}
	return &domain.Note{
		Id:         res.ID,
		UserId:     res.UserID,
		Title:      res.Title.String,
		Content:    res.Content.String,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		DeleteTime: res.DeleteTime.Time,
	}, nil
}

func (d *noteDatabaseDs) RestoreNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Note, error) {
	res, err := d.queries.WithTx(tx).RestoreNoteById(ctx, id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, err
		}
	}
	return &domain.Note{
		Id:         res.ID,
		UserId:     res.UserID,
		Title:      res.Title.String,
		Content:    res.Content.String,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		DeleteTime: res.DeleteTime.Time,
	}, nil
}

func (d *noteDatabaseDs) HardDeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	_, err := d.queries.WithTx(tx).HardDeleteNoteById(ctx, id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return &customerrors.RecordNotFound{}
		default:
			return err
		}
	}
	return nil
}

func (d *noteDatabaseDs) SoftDeleteNote(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	_, err := d.queries.WithTx(tx).SoftDeleteNoteById(ctx, database.SoftDeleteNoteByIdParams{
		ID:         id,
		DeleteTime: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return &customerrors.RecordNotFound{}
		default:
			return err
		}
	}
	return nil
}
