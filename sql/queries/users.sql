-- name: CreateUser :one
INSERT INTO users(email)
values(
  $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;
