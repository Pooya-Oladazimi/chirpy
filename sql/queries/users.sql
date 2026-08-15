-- name: CreateUser :one
INSERT INTO users(email)
values(
  $1
)
RETURNING *;

-- name: GetUser :one
SELECT * from users
where id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
