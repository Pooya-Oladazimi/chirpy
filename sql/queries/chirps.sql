-- name: CreateChirp :one
INSERT INTO chirps(body, user_id)
values(
  $1,
  $2
)
RETURNING *;

-- name: GetChirp :one
SELECT * from chirps
where id = $1;

-- name: GetAllChirps :many
SELECT * from chirps
order by created_at asc;

-- name: GetUserChirps :many
SELECT * from chirps
where user_id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
where id = $1;

