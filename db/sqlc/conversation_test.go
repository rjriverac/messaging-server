package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func createRandConv(t *testing.T) Conversation {

	name := util.RandomString(10)

	conv, err := testQueries.CreateConversation(context.Background(), sql.NullString{String: name, Valid: true})
	require.NoError(t, err)
	require.NotEmpty(t, conv)
	require.Equal(t, conv.Name.String, name)

	return conv
}

func TestCreateConv(t *testing.T) {
	createRandConv(t)
}

func TestGetConv(t *testing.T) {
	conv1 := createRandConv(t)

	conv2, err := testQueries.GetConversation(context.Background(), conv1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, conv2)
	require.Equal(t, conv1.ID, conv2.ID)
	require.Equal(t, conv1.Name, conv2.Name)
}

func TestListConv(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandConv(t)
	}

	arg := ListConversationsParams{
		Limit:  5,
		Offset: 5,
	}

	list, err := testQueries.ListConversations(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, list, 5)

	for _, conv := range list {
		require.NotEmpty(t, conv)
		require.Len(t, conv.Name.String, 10)
	}
}

func TestUpdateConv(t *testing.T) {
	conv1 := createRandConv(t)
	nName := util.RandomString(5)
	arg := UpdateConversationParams{ID: conv1.ID, Name: sql.NullString{String: nName, Valid: true}}

	conv2, err := testQueries.UpdateConversation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, conv2)
	require.Equal(t, conv2.Name.String, nName)
	require.Len(t, conv2.Name.String, 5)
}

func TestDeleteConv(t *testing.T) {
	conv1 := createRandConv(t)

	err := testQueries.DeleteConversation(context.Background(), conv1.ID)
	require.NoError(t, err)

	conv2, errnrow := testQueries.GetConversation(context.Background(), conv1.ID)
	require.Error(t, errnrow)
	require.Empty(t, conv2)
	require.EqualError(t, errnrow, sql.ErrNoRows.Error())
}

func TestGetConvMessages(t *testing.T) {
	conv1 := createRandConv(t)

	for i := 0; i < 20; i++ {
		arg := CreateMessageParams{
			From:    util.RandomString(10),
			Content: util.RandomString(50),
			ConvID:  conv1.ID,
		}
		_, err := testQueries.CreateMessage(context.Background(), arg)
		require.NoError(t, err)
	}
	list, err := testQueries.ListConvMessages(context.Background(), conv1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	now := time.Now()

	for _, msg := range list {

		require.NotEmpty(t, msg)
		require.Len(t, msg.From, 10)
		require.Len(t, msg.MessageContent, 50)
		require.NotZero(t, msg.MessageID)
		require.WithinDuration(t, now, msg.CreatedAt, time.Second)
	}
}
