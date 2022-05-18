-- name: CreateMessage :one
INSERT INTO "Message" ("user_id", "content")
VALUES($1, $2)
RETURNING *;
-- name: GetMessage :one
SELECT *
from "Message"
WHERE id = $1;
-- name: ListMessageByUser :many
SELECT *
from "Message"
WHERE "user_id" = $1
ORDER BY created_at;
-- name: DeleteMessage :exec
DELETE FROM "Message"
WHERE id = $1;