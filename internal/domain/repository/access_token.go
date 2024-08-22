package repository

import (
	"context"

	"github.com/daniarmas/notes/internal/entity"
	"github.com/google/uuid"
)

type AccessTokenRepository interface {
	GetAccessToken(ctx context.Context, id uuid.UUID) (*entity.AccessToken, error)
	CreateAccessToken(ctx context.Context, accessToken *entity.AccessToken) (*entity.AccessToken, error)
	DeleteAccessToken(ctx context.Context, id uuid.UUID) error
}
