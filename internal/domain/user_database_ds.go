package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserDatabaseDs interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
