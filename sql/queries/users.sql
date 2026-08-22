-- name: CreateUser :one
INSERT INTO users(email, hashed_password)
values(
  $1,
  $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUser :one
SELECT * from users
where id = $1;

-- name: GetUserByEmail :one
SELECT * from users
where email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
set email = sqlc.arg(new_email),
    hashed_password = sqlc.arg(password),
    updated_at = NOW()
WHERE id = sqlc.arg(user_id)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUserToChirpyRed :one
UPDATE users
set is_chirpy_red = TRUE
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: DowngradeUserFromChirpyRed :one
UPDATE users
set is_chirpy_red = FALSE
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;
