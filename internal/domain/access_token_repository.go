package domain

import (
	"context"
	"database/sql"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/utils"
	"github.com/google/uuid"
)

type AccessTokenRepository interface {
	GetAccessToken(ctx context.Context, id uuid.UUID) (*AccessToken, error)
	CreateAccessToken(ctx context.Context, tx *sql.Tx, userId uuid.UUID, refreshTokenId uuid.UUID) (*AccessToken, error)
	DeleteAccessTokenByUserId(ctx context.Context, tx *sql.Tx, userId uuid.UUID) error
}

type accessTokenRepository struct {
	AccessTokenCacheDs    AccessTokenCacheDs
	AccessTokenDatabaseDs AccessTokenDatabaseDs
}

func NewAccessTokenRepository(accessTokenCacheDs AccessTokenCacheDs, accessTokenDatabaseDs AccessTokenDatabaseDs) AccessTokenRepository {
	return &accessTokenRepository{
		AccessTokenCacheDs:    accessTokenCacheDs,
		AccessTokenDatabaseDs: accessTokenDatabaseDs,
	}
}

func (r *accessTokenRepository) GetAccessToken(ctx context.Context, id uuid.UUID) (*AccessToken, error) {
	// Get the access token from cache
	accessToken, err := r.AccessTokenCacheDs.GetAccessTokenById(ctx, id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			clogg.Error(
				ctx,
				"access token not found in cache",
				clogg.String("error", err.Error()),
				clogg.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
				clogg.String("file", utils.GetFileName()),
				clogg.String("function", utils.GetFunctionName()),
				clogg.Int("line", utils.GetLineNumber()),
			)
		default:
			clogg.Error(
				ctx,
				"error getting access token from cache",
				clogg.String("error", err.Error()),
				clogg.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
				clogg.String("file", utils.GetFileName()),
				clogg.String("function", utils.GetFunctionName()),
				clogg.Int("line", utils.GetLineNumber()),
			)
		}
		// Get the user from the database
		accessToken, err = r.AccessTokenDatabaseDs.GetAccessTokenById(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return accessToken, nil
}

func (r *accessTokenRepository) CreateAccessToken(ctx context.Context, tx *sql.Tx, userId uuid.UUID, refreshTokenId uuid.UUID) (*AccessToken, error) {
	accessToken := NewAccessToken(userId, refreshTokenId)
	// Save the access token in the database
	accessToken, err := r.AccessTokenDatabaseDs.CreateAccessToken(ctx, tx, accessToken)
	if err != nil {
		clogg.Error(
			ctx,
			"error creating access token in the database",
			clogg.String("error", err.Error()),
			clogg.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
			clogg.String("file", utils.GetFileName()),
			clogg.String("function", utils.GetFunctionName()),
			clogg.Int("line", utils.GetLineNumber()),
		)
		return nil, err
	}

	// Cache the access token
	err = r.AccessTokenCacheDs.CreateAccessToken(ctx, accessToken)
	if err != nil {
		clogg.Error(
			ctx,
			"error caching access token",
			clogg.String("error", err.Error()),
			clogg.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
			clogg.String("file", utils.GetFileName()),
			clogg.String("function", utils.GetFunctionName()),
			clogg.Int("line", utils.GetLineNumber()),
		)
	}

	return accessToken, nil
}

func (r *accessTokenRepository) DeleteAccessTokenByUserId(ctx context.Context, tx *sql.Tx, userId uuid.UUID) error {
	// Delete the refresh token on the database
	id, err := r.AccessTokenDatabaseDs.DeleteAccessTokenByUserId(ctx, tx, userId)
	if err != nil {
		return err
	}
	// Delete the refresh token on the cache
	err = r.AccessTokenCacheDs.DeleteAccessToken(ctx, *id)
	if err != nil {
		return err
	}
	return nil
}
