// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createAccessToken = `-- name: CreateAccessToken :one
INSERT INTO access_tokens (
  user_id, refresh_token_id
) VALUES (
  $1, $2
)
RETURNING id, user_id, refresh_token_id, create_time, update_time
`

type CreateAccessTokenParams struct {
	UserID         uuid.UUID
	RefreshTokenID uuid.UUID
}

func (q *Queries) CreateAccessToken(ctx context.Context, arg CreateAccessTokenParams) (AccessToken, error) {
	row := q.db.QueryRowContext(ctx, createAccessToken, arg.UserID, arg.RefreshTokenID)
	var i AccessToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshTokenID,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const createNote = `-- name: CreateNote :one
INSERT INTO notes (
  user_id, title, content, create_time, update_time
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, user_id, title, content, create_time, update_time, delete_time
`

type CreateNoteParams struct {
	UserID     uuid.UUID
	Title      sql.NullString
	Content    sql.NullString
	CreateTime time.Time
	UpdateTime time.Time
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.UserID,
		arg.Title,
		arg.Content,
		arg.CreateTime,
		arg.UpdateTime,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreateTime,
		&i.UpdateTime,
		&i.DeleteTime,
	)
	return i, err
}

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  user_id
) VALUES (
  $1
)
RETURNING id, user_id, create_time, update_time
`

func (q *Queries) CreateRefreshToken(ctx context.Context, userID uuid.UUID) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, userID)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING id, name, email, password, create_time, update_time
`

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const deleteAccessTokenByUserId = `-- name: DeleteAccessTokenByUserId :one
DELETE FROM access_tokens WHERE user_id = $1 RETURNING id
`

func (q *Queries) DeleteAccessTokenByUserId(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, deleteAccessTokenByUserId, userID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteRefreshTokenByUserId = `-- name: DeleteRefreshTokenByUserId :one
DELETE FROM refresh_tokens WHERE user_id = $1 RETURNING id
`

func (q *Queries) DeleteRefreshTokenByUserId(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, deleteRefreshTokenByUserId, userID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getAccessTokenById = `-- name: GetAccessTokenById :one
SELECT id, user_id, refresh_token_id, create_time, update_time FROM access_tokens
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccessTokenById(ctx context.Context, id uuid.UUID) (AccessToken, error) {
	row := q.db.QueryRowContext(ctx, getAccessTokenById, id)
	var i AccessToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshTokenID,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const getRefreshTokenById = `-- name: GetRefreshTokenById :one
SELECT id, user_id, create_time, update_time FROM refresh_tokens
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRefreshTokenById(ctx context.Context, id uuid.UUID) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshTokenById, id)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, create_time, update_time FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, name, email, password, create_time, update_time FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const hardDeleteNoteById = `-- name: HardDeleteNoteById :one
DELETE FROM notes WHERE id = $1 RETURNING id, user_id, title, content, create_time, update_time, delete_time
`

func (q *Queries) HardDeleteNoteById(ctx context.Context, id uuid.UUID) (Note, error) {
	row := q.db.QueryRowContext(ctx, hardDeleteNoteById, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreateTime,
		&i.UpdateTime,
		&i.DeleteTime,
	)
	return i, err
}

const listNotesByUserId = `-- name: ListNotesByUserId :many
SELECT id, user_id, title, content, create_time, update_time, delete_time FROM notes
WHERE user_id = $1 AND update_time < $2
ORDER BY update_time DESC
LIMIT 20
`

type ListNotesByUserIdParams struct {
	UserID     uuid.UUID
	UpdateTime time.Time
}

func (q *Queries) ListNotesByUserId(ctx context.Context, arg ListNotesByUserIdParams) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, listNotesByUserId, arg.UserID, arg.UpdateTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreateTime,
			&i.UpdateTime,
			&i.DeleteTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteNoteById = `-- name: SoftDeleteNoteById :one
UPDATE notes SET
  delete_time = $2
WHERE id = $1 AND delete_time IS NULL
RETURNING id, user_id, title, content, create_time, update_time, delete_time
`

type SoftDeleteNoteByIdParams struct {
	ID         uuid.UUID
	DeleteTime sql.NullTime
}

func (q *Queries) SoftDeleteNoteById(ctx context.Context, arg SoftDeleteNoteByIdParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, softDeleteNoteById, arg.ID, arg.DeleteTime)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreateTime,
		&i.UpdateTime,
		&i.DeleteTime,
	)
	return i, err
}

const updateNoteById = `-- name: UpdateNoteById :one
UPDATE notes SET
  title = $2, content = $3, update_time = $4
WHERE id = $1 RETURNING id, user_id, title, content, create_time, update_time, delete_time
`

type UpdateNoteByIdParams struct {
	ID         uuid.UUID
	Title      sql.NullString
	Content    sql.NullString
	UpdateTime time.Time
}

func (q *Queries) UpdateNoteById(ctx context.Context, arg UpdateNoteByIdParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, updateNoteById,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.UpdateTime,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreateTime,
		&i.UpdateTime,
		&i.DeleteTime,
	)
	return i, err
}
