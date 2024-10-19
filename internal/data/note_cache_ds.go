package data

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Note struct {
	Id              string    `redis:"id"`
	UserId          string    `redis:"user_id"`
	Title           string    `redis:"title"`
	Content         string    `redis:"content"`
	BackgroundColor string    `redis:"background_color"`
	CreateTime      time.Time `redis:"create_time"`
	UpdateTime      time.Time `redis:"update_time"`
	DeleteTime      time.Time `redis:"delete_time"`
}

// ParseToDomain converts a data.Note to a domain.Note
func (n *Note) parseToDomain() (*domain.Note, error) {
	// Check if the input note is nil
	if n == nil {
		return nil, nil
	}

	// Convert data.Note to domain.Note
	return &domain.Note{
		Id:              uuid.MustParse(n.Id),
		UserId:          uuid.MustParse(n.UserId),
		Title:           n.Title,
		Content:         n.Content,
		BackgroundColor: n.BackgroundColor,
		CreateTime:      n.CreateTime,
		UpdateTime:      n.UpdateTime,
		DeleteTime:      n.DeleteTime,
	}, nil
}

// parseNoteFromDomain converts a domain.Note to a data.Note
func parseFromDomain(note *domain.Note) *Note {
	// Check if the input note is nil
	if note == nil {
		return nil
	}

	// Convert domain.Note to data.Note
	return &Note{
		Id:              note.Id.String(),
		UserId:          note.UserId.String(),
		Title:           note.Title,
		Content:         note.Content,
		BackgroundColor: note.BackgroundColor,
		CreateTime:      note.CreateTime,
		UpdateTime:      note.UpdateTime,
		DeleteTime:      note.DeleteTime,
	}
}

type noteCacheDs struct {
	redis *redis.Client
}

func NewNoteCacheDs(redis *redis.Client) domain.NoteCacheDs {
	return &noteCacheDs{
		redis: redis,
	}
}

func (n *noteCacheDs) CreateNote(ctx context.Context, note *domain.Note) error {
	key := fmt.Sprintf("note:%s", note.Id)

	pipeline := n.redis.TxPipeline()

	// Add commands to the transaction
	pipeline.HSet(ctx, key, parseFromDomain(note)).Result()
	pipeline.Expire(ctx, key, 59*time.Minute) // Set expiration time

	// Execute the transaction
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
