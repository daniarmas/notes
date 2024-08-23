package domain

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenCacheDs interface {
	GetRefreshToken(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
}
