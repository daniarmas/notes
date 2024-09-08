package data

import (
	"context"
	"time"

	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type userDatabaseDs struct {
	queries *database.Queries
}

func New(queries *database.Queries) domain.UserDatabaseDs {
	return &userDatabaseDs{
		queries: queries,
	}
}

func (d *userDatabaseDs) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := d.queries.CreateUser(ctx, database.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Id:         res.ID,
		Name:       res.Name,
		Email:      res.Email,
		CreateTime: res.CreateTime,
	}, nil
}

func (d *userDatabaseDs) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	res, err := d.queries.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	// Check if time is not null
	var updateTime time.Time
	if res.UpdateTime.Valid {
		updateTime = res.UpdateTime.Time
	}
	return &domain.User{
		Id:         res.ID,
		Name:       res.Name,
		Password:   "",
		Email:      res.Email,
		CreateTime: res.CreateTime,
		UpdateTime: updateTime,
	}, nil
}