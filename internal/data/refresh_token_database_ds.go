package data

import (
	"context"
	"database/sql"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type refreshTokenDatabaseDs struct {
	queries *database.Queries
}

func parseRefreshTokenToDomain(refreshToken *database.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		Id:         refreshToken.ID,
		UserId:     refreshToken.UserID,
		CreateTime: refreshToken.CreateTime,
		UpdateTime: refreshToken.UpdateTime.Time,
	}
}

func NewRefreshTokenDatabaseDs(queries *database.Queries) domain.RefreshTokenDatabaseDs {
	return &refreshTokenDatabaseDs{
		queries: queries,
	}
}

func (d *refreshTokenDatabaseDs) CreateRefreshToken(ctx context.Context, tx *sql.Tx, refreshToken *domain.RefreshToken) (*domain.RefreshToken, error) {
	res, err := d.queries.WithTx(tx).CreateRefreshToken(ctx, refreshToken.UserId)
	if err != nil {
		switch err.Error() {
		case "ERROR: insert on table \"refresh_tokens\" violates foreign key constraint \"fk_user\" (SQLSTATE 23503)":
			return nil, &customerrors.ForeignKeyConstraint{Field: "user_id", ParentTable: "users"}
		default:
			return nil, err
		}
	}
	return parseRefreshTokenToDomain(&res), nil
}

func (d *refreshTokenDatabaseDs) GetRefreshTokenById(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	res, err := d.queries.GetRefreshTokenById(ctx, id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, err
		}
	}
	return parseRefreshTokenToDomain(&res), nil
}

func (d *refreshTokenDatabaseDs) DeleteRefreshTokenByUserId(ctx context.Context, tx *sql.Tx, userId uuid.UUID) (*uuid.UUID, error) {
	id, err := d.queries.WithTx(tx).DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, err
		}
	}
	return &id, nil
}
