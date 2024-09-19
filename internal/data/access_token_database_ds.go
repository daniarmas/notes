package data

import (
	"context"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

type accessTokenDatabaseDs struct {
	queries *database.Queries
}

func parseAccessTokenToDomain(accessToken *database.AccessToken) *domain.AccessToken {
	// Check if the input accessToken is nil
	if accessToken == nil {
		return nil
	}

	// Convert database.AccessToken to domain.AccessToken
	return &domain.AccessToken{
		Id:             accessToken.ID,
		UserId:         accessToken.UserID,
		RefreshTokenId: accessToken.RefreshTokenID,
		CreateTime:     accessToken.CreateTime,
		UpdateTime:     accessToken.UpdateTime.Time,
	}
}

func NewAccessTokenDatabaseDs(queries *database.Queries) domain.AccessTokenDatabaseDs {
	return &accessTokenDatabaseDs{
		queries: queries,
	}
}

func (d *accessTokenDatabaseDs) CreateAccessToken(ctx context.Context, accessToken *domain.AccessToken) (*domain.AccessToken, error) {
	res, err := d.queries.CreateAccessToken(ctx, database.CreateAccessTokenParams{
		UserID:         accessToken.UserId,
		RefreshTokenID: accessToken.RefreshTokenId,
	})
	if err != nil {
		switch err.Error() {
		case "ERROR: insert on table \"access_tokens\" violates foreign key constraint \"fk_refresh_token\" (SQLSTATE 23503)":
			return nil, &customerrors.ForeignKeyConstraint{Field: "refresh_token_id", ParentTable: "refresh_tokens"}
		case "ERROR: insert on table \"access_tokens\" violates foreign key constraint \"fk_user\" (SQLSTATE 23503)":
			return nil, &customerrors.ForeignKeyConstraint{Field: "user_id", ParentTable: "users"}
		default:
			return nil, err
		}
	}
	return parseAccessTokenToDomain(&res), nil
}

func (d *accessTokenDatabaseDs) GetAccessTokenId(ctx context.Context, id uuid.UUID) (*domain.AccessToken, error) {
	res, err := d.queries.GetAccessTokenById(ctx, id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, &customerrors.RecordNotFound{}
		default:
			return nil, err
		}
	}
	return parseAccessTokenToDomain(&res), nil
}

func (d *accessTokenDatabaseDs) DeleteAccessTokenByUserId(ctx context.Context, userId uuid.UUID) (*uuid.UUID, error) {
	id, err := d.queries.DeleteAccessTokenByUserId(ctx, userId)
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
