-- name: CreateUser :one
INSERT INTO "Users" (
    name,
    email,
    hashed_pw,
    image,
    status
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING $1,
  $2,
  $4,
  $5;
-- name: GetUser :one
SELECT id,
  name,
  email,
  image,
  status
FROM "Users"
WHERE id = $1
LIMIT 1;
-- name: ListUsers :many
SELECT id,
  name,
  email,
  image,
  status
FROM "Users"
ORDER BY id;