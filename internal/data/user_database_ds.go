package data

import (
	"context"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type userDatabaseDs struct {
	queries *database.Queries
}

func NewUserDatabaseDs(queries *database.Queries) domain.UserDatabaseDs {
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
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, &customerrors.DuplicateRecord{Field: "email"}
		default:
			return nil, &customerrors.Unknown{}
		}
	}
	return &domain.User{
		Id:         res.ID,
		Name:       res.Name,
		Email:      res.Email,
		Password:   res.Password,
		CreateTime: res.CreateTime,
	}, nil
}

func (d *userDatabaseDs) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	res, err := d.queries.GetUserById(ctx, id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, &customerrors.Unknown{}
		}
	}
	// Check if time is not null
	var updateTime time.Time
	if res.UpdateTime.Valid {
		updateTime = res.UpdateTime.Time
	}
	return &domain.User{
		Id:         res.ID,
		Name:       res.Name,
		Password:   res.Password,
		Email:      res.Email,
		CreateTime: res.CreateTime,
		UpdateTime: updateTime,
	}, nil
}

func (d *userDatabaseDs) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	res, err := d.queries.GetUserByEmail(ctx, email)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, &customerrors.Unknown{}
		}
	}
	// Check if time is not null
	var updateTime time.Time
	if res.UpdateTime.Valid {
		updateTime = res.UpdateTime.Time
	}
	return &domain.User{
		Id:         res.ID,
		Name:       res.Name,
		Password:   res.Password,
		Email:      res.Email,
		CreateTime: res.CreateTime,
		UpdateTime: updateTime,
	}, nil
}
