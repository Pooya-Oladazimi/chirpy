-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, user_id, expires_at, revoked_at)
values(
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * from refresh_tokens
where token = $1;



-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
set revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;

