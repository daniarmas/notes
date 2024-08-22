package database

import (
	"context"

	"github.com/daniarmas/notes/internal/entity"
)

type UserDatasource interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}
