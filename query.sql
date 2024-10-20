-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListNotesByUserId :many
SELECT * FROM notes
WHERE user_id = $1 AND create_time < $2
ORDER BY create_time DESC
LIMIT 20;

-- name: CreateNote :one
INSERT INTO notes (
  user_id, title, content, background_color
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateNoteById :one
UPDATE notes SET
  title = $2, content = $3, update_time = $4
WHERE id = $1 RETURNING *;

-- name: DeleteNoteById :one
DELETE FROM notes WHERE id = $1 RETURNING *;

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

-- name: DeleteAccessTokenByUserId :one
DELETE FROM access_tokens WHERE user_id = $1 RETURNING id;

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

-- name: DeleteRefreshTokenByUserId :one
DELETE FROM refresh_tokens WHERE user_id = $1 RETURNING id;