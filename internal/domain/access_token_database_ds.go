package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccessTokenDatabaseDs interface {
	GetAccessTokenId(ctx context.Context, id uuid.UUID) (*AccessToken, error)
	CreateAccessToken(ctx context.Context, accessToken *AccessToken) (*AccessToken, error)
	DeleteAccessTokenByUserId(ctx context.Context, userId uuid.UUID) (*uuid.UUID, error)
}
