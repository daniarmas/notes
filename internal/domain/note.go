package domain

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id           uuid.UUID `json:"id"`
	UserId       uuid.UUID `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	OriginalUrl  string    `json:"original_url"`
	ProcessedUrl string    `json:"processed_url"`
	Files        []*File   `json:"files"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	DeleteTime   time.Time `json:"delete_time"`
}
