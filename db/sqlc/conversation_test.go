package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func createRandConv(t *testing.T, uID1, uID2 int64) Conversation {

	msg1 := createRandMessage(t, uID1)
	msg2 := createRandMessage(t, uID2)

	arg := CreateConversationParams{
		Unread:   int32(util.RandomInt(1, 30)),
		Messages: msg1.ID,
		Last:     msg2.ID,
	}

	conv, err := testQueries.CreateConversation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, conv)

	require.Equal(t, arg.Last, conv.Last)
	require.Equal(t, arg.Messages, conv.Messages)
	require.Equal(t, arg.Unread, conv.Unread)

	return conv
}

func TestCreateConv(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	createRandConv(t, user1.ID, user2.ID)
}

func TestGetConv(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	conv1 := createRandConv(t, user1.ID,user2.ID)

	conv2, err := testQueries.GetConversation(context.Background(), conv1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, conv2)
	require.Equal(t, conv1.ID, conv2.ID)
	require.Equal(t, conv1.Last, conv2.Last)
	require.Equal(t, conv1.Messages, conv2.Messages)
	require.Equal(t, conv1.Unread, conv2.Unread)
}

func TestListConv(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandConv(t, user1.ID, user2.ID)
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
	}
}

func TestUpdateConv(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	conv1 := createRandConv(t, user1.ID, user2.ID)
	newMsgs := []Message{
		createRandMessage(t, user1.ID),
		createRandMessage(t, user1.ID),
	}
	arg := UpdateConversationParams{
		ID:       conv1.ID,
		Last:     newMsgs[0].ID,
		Unread:   int32(conv1.ID + 1),
		Messages: newMsgs[1].ID,
	}
	conv2,err := testQueries.UpdateConversation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t,conv2)
	require.Equal(t,arg.Last,conv2.Last)
	require.Equal(t,arg.Messages,conv2.Messages)
	require.Equal(t,arg.Unread,conv2.Unread)
}

func TestDeleteConv(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	conv1 := createRandConv(t,user1.ID,user2.ID)

	err := testQueries.DeleteConversation(context.Background(),conv1.ID)
	require.NoError(t,err)

	conv2, err := testQueries.GetConversation(context.Background(),conv1.ID)
	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,conv2)
}