-- name: CreateConversation :one
INSERT INTO "Conversation" (name)
VALUES($1)
RETURNING *;
-- name: GetConversation :one
SELECT *
FROM "Conversation"
WHERE id = $1
LIMIT 1;
-- name: ListConversations :many
SELECT *
FROM "Conversation"
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: UpdateConversation :one
UPDATE "Conversation"
SET name = $2
WHERE ID = $1
returning *;
-- name: DeleteConversation :exec
DELETE FROM "Conversation"
WHERE ID = $1;