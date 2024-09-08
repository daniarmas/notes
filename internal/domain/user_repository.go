package domain

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
}

type userRepo struct {
	UserCacheDs    UserCacheDs
	UserDatabaseDs UserDatabaseDs
}

func NewUserRepository(userCacheDs *UserCacheDs, userDatabaseDs *UserDatabaseDs) UserRepository {
	return &userRepo{
		UserCacheDs:    *userCacheDs,
		UserDatabaseDs: *userDatabaseDs,
	}
}

func (d *userRepo) CreateUser(ctx context.Context, user *User) (*User, error) {
	// Save the user on the database
	user, err := d.UserDatabaseDs.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	// Cache the user asynchronously, don't block the main operation
	go func() {
		err = d.UserCacheDs.CreateUser(ctx, user)
		if err != nil {
			log.Println(err)
		}
	}()
	return user, nil
}

func (d *userRepo) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	// Get the user from cache
	user, err := d.UserCacheDs.GetUser(ctx, id)
	if err != nil {
		log.Println(err)
		// Get the user from the database
		user, err = d.UserDatabaseDs.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
