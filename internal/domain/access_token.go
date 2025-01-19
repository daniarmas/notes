package domain

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

// NewAccessToken creates a new AccessToken instance
func NewAccessToken(userId uuid.UUID, refreshTokenId uuid.UUID) *AccessToken {
	now := time.Now()
	return &AccessToken{
		Id:             uuid.New(),
		UserId:         userId,
		RefreshTokenId: refreshTokenId,
		CreateTime:     now,
		UpdateTime:     now,
	}
}
