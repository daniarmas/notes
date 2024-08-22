package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	RefreshTokenId uuid.UUID `json:"refresh_token_id"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
}
