// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateConversation(ctx context.Context, name sql.NullString) (Conversation, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	CreateUser_conversation(ctx context.Context, arg CreateUser_conversationParams) (UserConversation, error)
	DeleteConversation(ctx context.Context, id int64) error
	DeleteMessage(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteUser_conversation(ctx context.Context, arg DeleteUser_conversationParams) error
	DeleteUser_conversation_by_id(ctx context.Context, id int64) error
	GetConversation(ctx context.Context, id int64) (Conversation, error)
	GetMessage(ctx context.Context, id int64) (Message, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	GetUser_conv_by_id(ctx context.Context, id int64) (UserConversation, error)
	GetUser_conversation(ctx context.Context, arg GetUser_conversationParams) (UserConversation, error)
	ListConversations(ctx context.Context, arg ListConversationsParams) ([]Conversation, error)
	ListMessageByUser(ctx context.Context, from string) ([]Message, error)
	ListUserMessages(ctx context.Context, id int64) ([]ListUserMessagesRow, error)
	ListUser_conversationByUser(ctx context.Context, userID int64) ([]UserConversation, error)
	ListUser_conversations(ctx context.Context) ([]UserConversation, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error)
	UpdateConversation(ctx context.Context, arg UpdateConversationParams) (Conversation, error)
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (UpdateUserInfoRow, error)
}

var _ Querier = (*Queries)(nil)
