package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserDatabaseDs interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
}
