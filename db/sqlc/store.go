package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store interface {
	Querier
	SendMessage(ctx context.Context, arg SendMessageParams) (SendResult, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

type SendMessageParams struct {
	Content string `json:"content"`
	ConvID  int64  `json:"convID"`
	UserID  int64  `json:"from_id"`
}

type SendResult struct {
	Timestamp time.Time `json:"sent_at"`
	MsgID     int64     `json:"id"`
}

func (store *SQLStore) SendMessage(ctx context.Context, arg SendMessageParams) (SendResult, error) {
	var result SendResult

	// validate conversation exists, create message
	// check join table

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		if _, err := q.GetUser_conversation(ctx, GetUser_conversationParams{
			UserID: arg.UserID,
			ConvID: arg.ConvID,
		}); err != nil {
			if err == sql.ErrNoRows {
				q.CreateUser_conversation(ctx, CreateUser_conversationParams{
					UserID: arg.UserID,
					ConvID: arg.ConvID,
				})
			}
		}
		user, err := q.GetUser(ctx, arg.UserID)
		if err != nil {
			return err
		}
		/*
			todo add query to db to get user and add
			todo the name to the message params here vs from the api side
		*/
		msg, err := q.CreateMessage(ctx, CreateMessageParams{
			From:    user.Name,
			Content: arg.Content,
			ConvID:  arg.ConvID,
		})
		if err != nil {
			return err
		}

		result.Timestamp = msg.CreatedAt
		result.MsgID = msg.ID

		return nil
	})
	return result, err
}
