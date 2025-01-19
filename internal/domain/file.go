package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id            uuid.UUID `json:"id"`
	NoteId        uuid.UUID `json:"note_id"`
	OriginalFile  string    `json:"original_file"`
	ProcessedFile string    `json:"processed_file"`
	Url           string    `json:"url"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	DeleteTime    time.Time `json:"delete_time"`
}
