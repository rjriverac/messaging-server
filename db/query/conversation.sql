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
-- name: ListConvMessages :many
SELECT 
"Message".from,"Message".content as message_content,"Message".created_at, "Message".id as message_id
FROM
"user_conversation"
INNER JOIN "Conversation" on "user_conversation".conv_id = "Conversation".id
INNER JOIN "Message" on "Conversation".id = "Message".conv_id
Where
"user_conversation".conv_id = $1
And "user_conversation".user_id=$2;