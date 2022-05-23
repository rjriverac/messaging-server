// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: conversation.sql

package db

import (
	"context"
	"database/sql"
)

const createConversation = `-- name: CreateConversation :one
INSERT INTO "Conversation" (name)
VALUES($1)
RETURNING id, name
`

func (q *Queries) CreateConversation(ctx context.Context, name sql.NullString) (Conversation, error) {
	row := q.db.QueryRowContext(ctx, createConversation, name)
	var i Conversation
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteConversation = `-- name: DeleteConversation :exec
DELETE FROM "Conversation"
WHERE ID = $1
`

func (q *Queries) DeleteConversation(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteConversation, id)
	return err
}

const getConversation = `-- name: GetConversation :one
SELECT id, name
FROM "Conversation"
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetConversation(ctx context.Context, id int64) (Conversation, error) {
	row := q.db.QueryRowContext(ctx, getConversation, id)
	var i Conversation
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listConversations = `-- name: ListConversations :many
SELECT id, name
FROM "Conversation"
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListConversationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListConversations(ctx context.Context, arg ListConversationsParams) ([]Conversation, error) {
	rows, err := q.db.QueryContext(ctx, listConversations, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conversation
	for rows.Next() {
		var i Conversation
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateConversation = `-- name: UpdateConversation :one
UPDATE "Conversation"
SET name = $2
WHERE ID = $1
returning id, name
`

type UpdateConversationParams struct {
	ID   int64          `json:"id"`
	Name sql.NullString `json:"name"`
}

func (q *Queries) UpdateConversation(ctx context.Context, arg UpdateConversationParams) (Conversation, error) {
	row := q.db.QueryRowContext(ctx, updateConversation, arg.ID, arg.Name)
	var i Conversation
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
