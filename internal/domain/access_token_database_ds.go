package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type AccessTokenDatabaseDs interface {
	GetAccessTokenById(ctx context.Context, id uuid.UUID) (*AccessToken, error)
	CreateAccessToken(ctx context.Context, tx *sql.Tx, accessToken *AccessToken) (*AccessToken, error)
	DeleteAccessTokenByUserId(ctx context.Context, tx *sql.Tx, userId uuid.UUID) (*uuid.UUID, error)
}
