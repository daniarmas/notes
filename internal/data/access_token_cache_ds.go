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

type AccessToken struct {
	Id             string    `redis:"id"`
	RefreshTokenId string    `redis:"refresh_token_id"`
	UserId         string    `redis:"user_id"`
	CreateTime     time.Time `redis:"create_time"`
	UpdateTime     time.Time `redis:"update_time"`
}

func (u *AccessToken) ParseToDomain() *domain.AccessToken {
	if u.Id != "" {
		id := uuid.MustParse(u.Id)
		userId := uuid.MustParse(u.UserId)
		refreshTokenId := uuid.MustParse(u.RefreshTokenId)
		return &domain.AccessToken{
			Id:             id,
			UserId:         userId,
			RefreshTokenId: refreshTokenId,
			CreateTime:     u.CreateTime,
			UpdateTime:     u.UpdateTime,
		}
	}
	return nil
}

func parseAccessTokenFromDomain(accessToken *domain.AccessToken) *AccessToken {
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

func (ds *accessTokenCacheDs) GetAccessToken(ctx context.Context, id uuid.UUID) (*domain.AccessToken, error) {
	key := fmt.Sprintf("access_token:%s", id)
	var response AccessToken
	if err := ds.redis.HGetAll(ctx, key).Scan(&response); err != nil {
		return nil, &customerrors.Unknown{}
	}
	if response.Id == "" {
		return nil, &customerrors.RecordNotFound{}
	}
	return response.ParseToDomain(), nil
}

func (ds *accessTokenCacheDs) CreateAccessToken(ctx context.Context, accessToken *domain.AccessToken) error {
	key := fmt.Sprintf("access_token:%s", accessToken.Id)
	_, err := ds.redis.HSet(ctx, key, parseAccessTokenFromDomain(accessToken)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ds *accessTokenCacheDs) DeleteAccessToken(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("access_token:%s", id)
	if _, err := ds.redis.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
