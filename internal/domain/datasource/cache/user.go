package database

import (
	"context"

	"github.com/daniarmas/notes/internal/entity"
	"github.com/google/uuid"
)

type UserCacheDs interface {
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
