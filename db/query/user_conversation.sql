-- name: CreateUser_conversation :one
INSERT INTO "user_conversation" (UserID, ConvID)
VALUES($1, $2)
RETURNING *;
-- name: GetUser_conversation :one
SELECT *
from "user_conversation"
WHERE user_conversation_pkey = $1;
-- name: ListUser_conversation :many
SELECT *
from "user_conversation"
ORDER BY UserID;
-- name: DeleteUser_conversation :exec
DELETE FROM "user_conversation"
WHERE user_conversation_pkey = $1;