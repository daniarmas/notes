package repository

import (
	"context"

	"github.com/daniarmas/notes/internal/entity"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, refreshToken *entity.RefreshToken) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
}