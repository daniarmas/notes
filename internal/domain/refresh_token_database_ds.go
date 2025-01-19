package domain

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenDatabaseDs interface {
	GetRefreshTokenById(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) (*RefreshToken, error)
	DeleteRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) (*uuid.UUID, error)
}
