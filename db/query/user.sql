-- name: CreateUser :one
INSERT INTO "Users" (
    name,
    email,
    hashed_pw,
    image,
    status
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
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
-- name: UpdateUserInfo :one
UPDATE "Users"
SET 
    name = coalesce(sqlc.narg('name'), name),
    email = coalesce(sqlc.narg('email'), email),
    image = coalesce(sqlc.narg('image'), image),
    status = coalesce(sqlc.narg('status'), status),
    hashed_pw = coalesce(sqlc.narg('hashed_pw'), hashed_pw)
where id = sqlc.arg('id')
RETURNING id,
  name,
  email,
  image,
  status,
  created_at;
-- name: DeleteUser :exec
DELETE FROM "Users"
WHERE id = $1;
-- name: ListUserMessages :many
SELECT 
"Conversation".name as conversation_name,"Message".from,"Message".content as message_content,"Message".created_at,
"user_conversation".conv_id, "Message".id as message_id
FROM
"Users"
INNER JOIN "user_conversation" on "Users".id = "user_conversation".user_id
INNER JOIN "Conversation" on "user_conversation".conv_id = "Conversation".id
inner JOIN "Message" on "Conversation".id = "Message".conv_id
Where "Users".id = $1;
-- name: ListConvFromUser :many
SELECT 
"Conversation".id,"Conversation".name
FROM
"Users"
INNER JOIN "user_conversation" on "Users".id = "user_conversation".user_id
INNER JOIN "Conversation" on "user_conversation".conv_id = "Conversation".id
WHERE "Users".id = $1;