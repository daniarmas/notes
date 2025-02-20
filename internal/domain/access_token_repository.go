package domain

import (
	"context"
	"database/sql"
	"log/slog"

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
			slog.LogAttrs(
				context.Background(),
				slog.LevelInfo,
				"Access token not found in cache",
				slog.String("error", err.Error()),
				slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
				slog.String("file", utils.GetFileName()),
				slog.String("function", utils.GetFunctionName()),
				slog.Int("line", utils.GetLineNumber()),
			)
		default:
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"Error getting access token from cache",
				slog.String("error", err.Error()),
				slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
				slog.String("file", utils.GetFileName()),
				slog.String("function", utils.GetFunctionName()),
				slog.Int("line", utils.GetLineNumber()),
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
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"Error creating access token in the database",
			slog.String("error", err.Error()),
			slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
			slog.String("file", utils.GetFileName()),
			slog.String("function", utils.GetFunctionName()),
			slog.Int("line", utils.GetLineNumber()),
		)
		return nil, err
	}

	// Cache the access token
	err = r.AccessTokenCacheDs.CreateAccessToken(ctx, accessToken)
	if err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"Error caching access token",
			slog.String("error", err.Error()),
			slog.String("request_id", utils.ExtractRequestIdFromContext(ctx)),
			slog.String("file", utils.GetFileName()),
			slog.String("function", utils.GetFunctionName()),
			slog.Int("line", utils.GetLineNumber()),
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
