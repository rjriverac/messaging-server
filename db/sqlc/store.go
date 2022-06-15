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
	CreateConv(ctx context.Context, arg CreateConvParams) (ConvReturn, error)
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

type NString string

func (s *NString) toNullStr() sql.NullString {
	if len(*s) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

type NullString sql.NullString

func (ns NullString) MarshalJson() string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

type CreateConvParams struct {
	Name    NString  `json:"name"`
	ToUsers []string `json:"recipient_emails"`
	From    int64    `json:"from"`
}

type ConvReturn struct {
	Name string `json:"conv_name"`
	ID   int64  `json:"conv_id"`
}

func (store *SQLStore) CreateConv(ctx context.Context, convParams CreateConvParams) (ConvReturn, error) {
	var ret ConvReturn

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		var validUsers []User

		for _, email := range convParams.ToUsers {
			user, err := q.GetUserByEmail(ctx, email)
			if err != nil {
				continue
			}
			validUsers = append(validUsers, user)
		}
		if len(validUsers) == 0 {
			return err
		}

		conv, err := q.CreateConversation(ctx, convParams.Name.toNullStr())
		if err != nil {
			return err
		}
		q.CreateUser_conversation(ctx, CreateUser_conversationParams{UserID: convParams.From, ConvID: conv.ID})

		for _, user := range validUsers {
			if _, err := q.GetUser_conversation(ctx, GetUser_conversationParams{UserID: user.ID, ConvID: conv.ID}); err != nil {
				if err == sql.ErrNoRows {
					q.CreateUser_conversation(ctx, CreateUser_conversationParams{UserID: user.ID, ConvID: conv.ID})
				}
				continue
			}
		}
		if err != nil {
			return err
		}
		ret.Name = NullString(conv.Name).MarshalJson()
		ret.ID = conv.ID
		return nil
	})
	return ret, err
}
