package domain

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) (*RefreshToken, error)
	DeleteRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) error
}

type refreshTokenRepository struct {
	RefreshTokenCacheDs    RefreshTokenCacheDs
	RefreshTokenDatabaseDs RefreshTokenDatabaseDs
}

func NewRefreshTokenRepository(refreshTokenCacheDs *RefreshTokenCacheDs, refreshTokenDatabaseDs *RefreshTokenDatabaseDs) RefreshTokenRepository {
	return &refreshTokenRepository{
		RefreshTokenCacheDs:    *refreshTokenCacheDs,
		RefreshTokenDatabaseDs: *refreshTokenDatabaseDs,
	}
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, id uuid.UUID) (*RefreshToken, error) {
	// Get the refresh token from cache
	refreshToken, err := r.RefreshTokenCacheDs.GetRefreshToken(ctx, id)
	if err != nil {
		log.Println(err)
		// Get the refresh from the database
		refreshToken, err = r.RefreshTokenDatabaseDs.GetRefreshTokenById(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return refreshToken, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) (*RefreshToken, error) {
	// Save the refresh token on the database
	user, err := r.RefreshTokenDatabaseDs.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	// Cache the refresh token asynchronously, don't block the main operation
	go func() {
		err = r.RefreshTokenCacheDs.CreateRefreshToken(ctx, user)
		if err != nil {
			log.Println(err)
		}
	}()
	return user, nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, userId uuid.UUID) error {
	// Delete the refresh token on the database
	id, err := r.RefreshTokenDatabaseDs.DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		return err
	}
	// Delete the refresh token on the cache
	err = r.RefreshTokenCacheDs.DeleteRefreshToken(ctx, *id)
	if err != nil {
		return err
	}
	return nil
}
