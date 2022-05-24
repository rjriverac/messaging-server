-- name: CreateUser_conversation :one
INSERT INTO "user_conversation" (user_id, conv_id)
VALUES($1, $2)
RETURNING *;
-- name: GetUser_conversation :one
SELECT *
from "user_conversation"
WHERE user_id = $1
  and conv_id = $2;
-- name: GetUser_conv_by_id :one
SELECT *
from "user_conversation"
WHERE id = $1;
-- name: ListUser_conversationByUser :many
SELECT *
from "user_conversation"
WHERE user_id = $1
ORDER BY user_id;
-- name: ListUser_conversations :many
SELECT *
from "user_conversation"
ORDER BY id;
-- name: DeleteUser_conversation :exec
DELETE FROM "user_conversation"
WHERE user_id = $1
  and conv_id = $2;
-- name: DeleteUser_conversation_by_id :exec
DELETE FROM "user_conversation"
WHERE id = $1;