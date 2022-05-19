-- name: CreateUser_conversation :one
INSERT INTO "user_conversation" (user_id, conv_id)
VALUES($1, $2)
RETURNING *;
-- name: GetUser_conversation :one
SELECT *
from "user_conversation"
WHERE user_id = $1 and conv_id = $2;
-- name: ListUser_conversationByUser :many
SELECT *
from "user_conversation"
WHERE user_id = $1
ORDER BY user_id;
-- name: ListUser_conversationByConv :many
SELECT *
from "user_conversation"
WHERE conv_id = $1
ORDER BY conv_id;
-- name: DeleteUser_conversation :exec
DELETE FROM "user_conversation"
WHERE user_id = $1 and conv_id = $2;