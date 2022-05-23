-- name: CreateMessage :one
INSERT INTO "Message" ("from", content, conv_id)
VALUES($1, $2, $3)
RETURNING *;
-- name: GetMessage :one
SELECT *
from "Message"
WHERE id = $1;
-- name: ListMessageByUser :many
SELECT *
from "Message"
WHERE "from" = $1
ORDER BY created_at;
-- name: DeleteMessage :exec
DELETE FROM "Message"
WHERE id = $1;