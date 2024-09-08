package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccessTokenCacheDs interface {
	GetAccessToken(ctx context.Context, id uuid.UUID) (*AccessToken, error)
	CreateAccessToken(ctx context.Context, accessToken *AccessToken) error
	DeleteAccessToken(ctx context.Context, id uuid.UUID) error
}
