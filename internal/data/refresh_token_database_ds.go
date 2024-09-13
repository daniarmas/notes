package data

import (
	"context"

	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type refreshTokenDatabaseDs struct {
	queries *database.Queries
}

func parseRefreshTokenToDomain(refreshToken *database.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		Id:         refreshToken.ID,
		UserId:     refreshToken.UserID,
		CreateTime: refreshToken.CreateTime,
		UpdateTime: refreshToken.UpdateTime.Time,
	}
}

func NewRefreshTokenDatabaseDs(queries *database.Queries) domain.RefreshTokenDatabaseDs {
	return &refreshTokenDatabaseDs{
		queries: queries,
	}
}

func (d *refreshTokenDatabaseDs) CreateRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) (*domain.RefreshToken, error) {
	res, err := d.queries.CreateRefreshToken(ctx, refreshToken.UserId)
	if err != nil {
		return nil, err
	}
	return parseRefreshTokenToDomain(&res), nil
}

func (d *refreshTokenDatabaseDs) GetRefreshTokenById(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	res, err := d.queries.GetRefreshTokenById(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseRefreshTokenToDomain(&res), nil
}

func (d *refreshTokenDatabaseDs) DeleteRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) (*uuid.UUID, error) {
	id, err := d.queries.DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
