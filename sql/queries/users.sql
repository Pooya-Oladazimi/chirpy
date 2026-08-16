-- name: CreateUser :one
INSERT INTO users(email, hashed_password)
values(
  $1,
  $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetUser :one
SELECT * from users
where id = $1;

-- name: GetUserByEmail :one
SELECT * from users
where email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
