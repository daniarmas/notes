package domaindatabase

import (
	"context"

	"github.com/daniarmas/notes/internal/domain/entity"
	"github.com/google/uuid"
)

type UserDatabaseDs interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
