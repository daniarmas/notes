package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id          uuid.UUID `json:"id"`
	NoteId      uuid.UUID `json:"note_id"`
	OriginalId  string    `json:"original_id"`
	ProcessedId string    `json:"processed_id"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	DeleteTime  time.Time `json:"delete_time"`
}
