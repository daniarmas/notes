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

-- name: CreateNote :one
INSERT INTO notes (
  user_id, title, content, background_color
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;