-- name: CreateUser :one
INSERT INTO "Users" (
    name,
    email,
    hashed_pw,
    image,
    status
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING id,
  name,
  email,
  image,
  status,
  created_at;
-- name: GetUser :one
SELECT id,
  name,
  email,
  image,
  status,
  created_at
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
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: UpdateUserImage :one
UPDATE "Users"
SET image = $2
WHERE id = $1
RETURNING id,
  name,
  email,
  image,
  status,
  created_at;
-- name: UpdateUserInfo :one
UPDATE "Users"
SET (name, email, image, status, hashed_pw) = ($2, $3, $4, $5, $6)
where id = $1
RETURNING id,
  name,
  email,
  image,
  status,
  created_at;
-- name: DeleteUser :exec
DELETE FROM "Users"
WHERE id = $1;
-- -- name: UpdateUserHashedPW :exec
-- UPDATE "Users"
-- SET hashed_pw = $2
-- WHERE id = $1;
-- -- name: UpdateUserEmail :exec
-- UPDATE "Users"
-- SET email = $2
-- WHERE id = $1;
-- -- name: UpdateUsername :exec
-- UPDATE "Users"
-- SET name = $2
-- WHERE id = $1;