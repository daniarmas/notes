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
WHERE user_id = $1 AND update_time < $2 AND delete_time IS NULL
ORDER BY update_time DESC
LIMIT 20;

-- name: ListTrashNotesByUserId :many
SELECT * FROM notes
WHERE user_id = $1 AND delete_time < $2 AND delete_time IS NOT NULL
ORDER BY delete_time DESC
LIMIT 20;

-- name: CreateNote :one
INSERT INTO notes (
  user_id, title, content, create_time, update_time
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateNoteById :one
UPDATE notes SET
  title = $2, content = $3, update_time = $4
WHERE id = $1 RETURNING *;

-- name: RestoreNoteById :one
UPDATE notes SET
  delete_time = NULL
WHERE id = $1 AND delete_time IS NOT NULL
RETURNING *;

-- name: HardDeleteNoteById :one
DELETE FROM notes WHERE id = $1 RETURNING *;

-- name: SoftDeleteNoteById :one
UPDATE notes SET
  delete_time = $2
WHERE id = $1 AND delete_time IS NULL
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

-- name: CreateFile :one
INSERT INTO files (
  note_id, original_file, create_time, update_time
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateFileByOriginalId :one
UPDATE files SET
  processed_file = $2, update_time = $3
WHERE original_file = $1 AND processed_file == '' RETURNING *;