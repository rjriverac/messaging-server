// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID   int64          `json:"id"`
	Name sql.NullString `json:"name"`
}

type Message struct {
	ID        int64         `json:"id"`
	From      string        `json:"from"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"createdAt"`
	ConvID    sql.NullInt64 `json:"convID"`
}

type User struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	HashedPw  string         `json:"hashedPw"`
	Image     sql.NullString `json:"image"`
	Status    sql.NullString `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
}

type UserConversation struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userID"`
	ConvID int64 `json:"convID"`
}
