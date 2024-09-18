package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AccessToken struct {
	Id             string    `redis:"id"`
	RefreshTokenId string    `redis:"refresh_token_id"`
	UserId         string    `redis:"user_id"`
	CreateTime     time.Time `redis:"create_time"`
	UpdateTime     time.Time `redis:"update_time"`
}

func (u *AccessToken) ParseToDomain() (*domain.AccessToken, error) {
	if u.Id == "" {
		return nil, nil
	}

	// Parse UUIDs and handle potential errors
	id, err := uuid.Parse(u.Id)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(u.UserId)
	if err != nil {
		return nil, err
	}
	refreshTokenId, err := uuid.Parse(u.RefreshTokenId)
	if err != nil {
		return nil, err
	}

	// Return the parsed domain.AccessToken
	return &domain.AccessToken{
		Id:             id,
		UserId:         userId,
		RefreshTokenId: refreshTokenId,
		CreateTime:     u.CreateTime,
		UpdateTime:     u.UpdateTime,
	}, nil
}

func parseAccessTokenFromDomain(accessToken *domain.AccessToken) *AccessToken {
	// Check if the input accessToken is nil
	if accessToken == nil {
		return nil
	}

	// Convert domain.AccessToken to AccessToken
	return &AccessToken{
		Id:             accessToken.Id.String(),
		UserId:         accessToken.UserId.String(),
		RefreshTokenId: accessToken.RefreshTokenId.String(),
		CreateTime:     accessToken.CreateTime,
		UpdateTime:     accessToken.UpdateTime,
	}
}

type accessTokenCacheDs struct {
	redis *redis.Client
}

func NewAccessTokenTokenCacheDs(redis *redis.Client) domain.AccessTokenCacheDs {
	return &accessTokenCacheDs{
		redis: redis,
	}
}

func (ds *accessTokenCacheDs) GetAccessTokenById(ctx context.Context, id uuid.UUID) (*domain.AccessToken, error) {
	key := fmt.Sprintf("access_token:%s", id)
	var response AccessToken
	if err := ds.redis.HGetAll(ctx, key).Scan(&response); err != nil {
		return nil, &customerrors.Unknown{}
	}
	if response.Id == "" {
		return nil, &customerrors.RecordNotFound{}
	}
	return response.ParseToDomain()
}

func (ds *accessTokenCacheDs) CreateAccessToken(ctx context.Context, accessToken *domain.AccessToken) error {
	key := fmt.Sprintf("access_token:%s", accessToken.Id)

	pipeline := ds.redis.TxPipeline()

	// Add commands to the transaction
	pipeline.HSet(ctx, key, parseAccessTokenFromDomain(accessToken)).Result()
	pipeline.Expire(ctx, key, 59*time.Minute) // Set expiration time

	// Execute the transaction
	_, err := pipeline.Exec(ctx)
	if err != nil {
		slog.Error(err.Error())
		return &customerrors.Unknown{}
	}
	return nil
}

func (ds *accessTokenCacheDs) DeleteAccessToken(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("access_token:%s", id)
	if _, err := ds.redis.Del(ctx, key).Result(); err != nil {
		return &customerrors.Unknown{}
	}
	return nil
}
