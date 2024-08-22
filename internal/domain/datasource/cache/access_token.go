package domaincache

import (
	"context"

	"github.com/daniarmas/notes/internal/domain/entity"
	"github.com/google/uuid"
)

type AccessTokenCacheDs interface {
	GetAccessToken(ctx context.Context, id uuid.UUID) (*entity.AccessToken, error)
	CreateAccessToken(ctx context.Context, accessToken *entity.AccessToken) (*entity.AccessToken, error)
	DeleteAccessToken(ctx context.Context, id uuid.UUID) error
}
