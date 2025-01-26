package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RefreshTokenDatabaseDs interface {
	GetRefreshTokenById(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	CreateRefreshToken(ctx context.Context, tx *sql.Tx, refreshToken *RefreshToken) (*RefreshToken, error)
	DeleteRefreshTokenByUserId(ctx context.Context, tx *sql.Tx, userId uuid.UUID) (*uuid.UUID, error)
}
