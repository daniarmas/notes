-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListNotesByUserId :many
SELECT * FROM notes WHERE user_id = $1;

-- name: CreateNote :one
INSERT INTO notes (
  user_id, title, content, background_color
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccessTokenById :one
SELECT * FROM access_tokens
WHERE id = $1 LIMIT 1;

-- name: CreateAccessToken :one
INSERT INTO access_tokens (
  user_id, refresh_token_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteAccessTokenById :exec
DELETE FROM access_tokens WHERE id = $1;

-- name: GetRefreshTokenById :one
SELECT * FROM refresh_tokens
WHERE id = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  user_id
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteRefreshTokenById :exec
DELETE FROM refresh_tokens WHERE id = $1;