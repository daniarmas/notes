package data

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshToken struct {
	Id         string    `redis:"id"`
	UserId     string    `redis:"user_id"`
	CreateTime time.Time `redis:"create_time"`
	UpdateTime time.Time `redis:"update_time"`
}

func (u *RefreshToken) ParseToDomain() *domain.RefreshToken {
	id := uuid.MustParse(u.Id)
	userId := uuid.MustParse(u.UserId)
	return &domain.RefreshToken{
		Id:         id,
		UserId:     userId,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

func parseRefreshTokenFromDomain(refreshToken *domain.RefreshToken) *RefreshToken {
	return &RefreshToken{
		Id:         refreshToken.Id.String(),
		UserId:     refreshToken.UserId.String(),
		CreateTime: refreshToken.CreateTime,
		UpdateTime: refreshToken.UpdateTime,
	}
}

type refreshTokenCacheDs struct {
	redis *redis.Client
}

func NewRefreshTokenCacheDs(redis *redis.Client) domain.RefreshTokenCacheDs {
	return &refreshTokenCacheDs{
		redis: redis,
	}
}

func (ds *refreshTokenCacheDs) GetRefreshToken(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	key := fmt.Sprintf("refresh_token:%s", id)
	var response RefreshToken
	if err := ds.redis.HGetAll(ctx, key).Scan(&response); err != nil {
		return nil, err
	}
	return response.ParseToDomain(), nil
}

func (ds *refreshTokenCacheDs) CreateRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken.Id)
	_, err := ds.redis.HSet(ctx, key, parseRefreshTokenFromDomain(refreshToken)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ds *refreshTokenCacheDs) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("refresh_token:%s", id)
	if _, err := ds.redis.Del(ctx, key).Result(); err != nil {
		return &customerrors.Unknown{}
	}
	return nil
}
