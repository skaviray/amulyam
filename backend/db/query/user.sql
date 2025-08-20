-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  email,
  full_name,
  role,
  is_admin
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2,
    password_changed_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
