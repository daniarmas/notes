package data

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Id         string    `redis:"id"`
	Name       string    `redis:"name"`
	Email      string    `redis:"email"`
	Password   string    `redis:"password"`
	CreateTime time.Time `redis:"create_time"`
	UpdateTime time.Time `redis:"update_time"`
}

func (u *User) ParseToDomain() *domain.User {
	id := uuid.MustParse(u.Id)
	return &domain.User{
		Id:         id,
		Name:       u.Name,
		Email:      u.Email,
		Password:   u.Password,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

func parseUserFromDomain(user *domain.User) *User {
	return &User{
		Id:         user.Id.String(),
		Name:       user.Name,
		Email:      user.Email,
		Password:   user.Password,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}

type userCacheDs struct {
	redis *redis.Client
}

func NewUserCacheDs(redis *redis.Client) domain.UserCacheDs {
	return &userCacheDs{
		redis: redis,
	}
}

func (ds *userCacheDs) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", id)
	var response User
	if err := ds.redis.HGetAll(ctx, key).Scan(&response); err != nil {
		return nil, err
	}
	return response.ParseToDomain(), nil
}

func (ds *userCacheDs) CreateUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.Id)
	_, err := ds.redis.HSet(ctx, key, parseUserFromDomain(user)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ds *userCacheDs) UpdateUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.Id)
	_, err := ds.redis.HSet(ctx, key, parseUserFromDomain(user)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ds *userCacheDs) DeleteUser(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("user:%s", id)
	if _, err := ds.redis.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
