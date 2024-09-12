package domain

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type AccessTokenRepository interface {
	GetAccessToken(ctx context.Context, id uuid.UUID) (*AccessToken, error)
	CreateAccessToken(ctx context.Context, accessToken *AccessToken) (*AccessToken, error)
	DeleteAccessToken(ctx context.Context, id uuid.UUID) error
}

type accessTokenRepository struct {
	AccessTokenCacheDs    AccessTokenCacheDs
	AccessTokenDatabaseDs AccessTokenDatabaseDs
}

func NewAccessTokenRepository(accessTokenCacheDs *AccessTokenCacheDs, accessTokenDatabaseDs *AccessTokenDatabaseDs) AccessTokenRepository {
	return &accessTokenRepository{
		AccessTokenCacheDs:    *accessTokenCacheDs,
		AccessTokenDatabaseDs: *accessTokenDatabaseDs,
	}
}

func (r *accessTokenRepository) GetAccessToken(ctx context.Context, id uuid.UUID) (*AccessToken, error) {
	// Get the access token from cache
	accessToken, err := r.AccessTokenCacheDs.GetAccessToken(ctx, id)
	if err != nil {
		log.Println(err)
		// Get the user from the database
		accessToken, err = r.AccessTokenDatabaseDs.GetAccessToken(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return accessToken, nil
}

func (r *accessTokenRepository) CreateAccessToken(ctx context.Context, accessToken *AccessToken) (*AccessToken, error) {
	// Save the access token on the database
	user, err := r.AccessTokenDatabaseDs.CreateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// Cache the access token asynchronously, don't block the main operation
	go func() {
		err = r.AccessTokenCacheDs.CreateAccessToken(ctx, user)
		if err != nil {
			log.Println(err)
		}
	}()
	return user, nil
}

func (r *accessTokenRepository) DeleteAccessToken(ctx context.Context, id uuid.UUID) error {
	// Delete the access token on the database
	err := r.AccessTokenDatabaseDs.DeleteAccessToken(ctx, id)
	if err != nil {
		return err
	}
	err = r.AccessTokenCacheDs.DeleteAccessToken(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
