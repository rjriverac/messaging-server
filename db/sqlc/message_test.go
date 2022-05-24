package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func createRandMessage(t *testing.T) Message {

	user := createRandomUser(t)

	conv := createRandConv(t)

	arg := CreateMessageParams{
		Content: util.RandomString(50),
		ConvID: sql.NullInt64{conv.ID,true},
		From: user.Name,
	}
	message, err := testQueries.CreateMessage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message)

	require.Equal(t, arg.Content, message.Content)
	require.Equal(t, arg.From, message.From)
	require.NotEmpty(t, message.CreatedAt)

	return message
}

func TestCreateMessage(t *testing.T) {
	createRandMessage(t)
}

func TestGetMessage(t *testing.T) {

	msg1 := createRandMessage(t)
	msg2, err := testQueries.GetMessage(context.Background(), msg1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, msg2)

	require.Equal(t, msg1.Content, msg2.Content)
	require.Equal(t, msg1.ID, msg2.ID)
	require.Equal(t, msg1.From, msg2.From)
	require.WithinDuration(t, msg1.CreatedAt, msg2.CreatedAt, time.Second)
}

func TestListMessages(t *testing.T) {

	user := createRandomUser(t)
	conv := createRandConv(t)
	arg := CreateMessageParams{
		Content: util.RandomString(50),
		ConvID: sql.NullInt64{conv.ID,true},
		From: user.Name,
	}

	for i := 0; i < 20; i++ {
		msg,err := testQueries.CreateMessage(context.Background(),arg)
		require.NoError(t,err)
		require.NotEmpty(t,msg)
	}

	messages, err := testQueries.ListMessageByUser(context.Background(), user.Name)

	require.NoError(t, err)
	require.Len(t, messages, 20)

	for _, msg := range messages {
		require.NotEmpty(t, msg)
	}

}

func TestDeleteMessage(t *testing.T) {
	

	msg := createRandMessage(t)
	err := testQueries.DeleteMessage(context.Background(),msg.ID)
	require.NoError(t,err)

	msg2, err := testQueries.GetMessage(context.Background(), msg.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, msg2)
}
