-- name: CreateConversation :one
INSERT INTO "Conversation" (Unread, Last, Messages)
VALUES ($1, $2, $3)
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
-- name: UpdateConversation :exec
UPDATE "Conversation"
SET (Unread, Last, Messages) = ($2, $3, $4)
WHERE ID = $1
returning *;
-- name: DeleteConversation :exec
DELETE FROM "Conversation"
WHERE ID = $1;