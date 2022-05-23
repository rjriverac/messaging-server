package db

// import (
// 	"context"

// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func createRndUsrConv(t *testing.T, uID int64, cID int64) UserConversation {

// 	arg := CreateUser_conversationParams{
// 		UserID: uID,
// 		ConvID: cID,
// 	}
// 	join, err := testQueries.CreateUser_conversation(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, join)
// 	require.Equal(t, arg.ConvID, join.ConvID)
// 	require.Equal(t, arg.UserID, join.UserID)

// 	return join
// }

// func TestCreateUsrConv(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	user2 := createRandomUser(t)
// 	conv1 := createRandConv(t, user1.ID, user2.ID)
// 	createRndUsrConv(t, user1.ID, conv1.ID)
// }

// func TestGetUserConv(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	user2 := createRandomUser(t)
// 	conv1 := createRandConv(t, user1.ID, user2.ID)
// 	uconv := createRndUsrConv(t, user1.ID, conv1.ID)

// 	arg := GetUser_conversationParams{
// 		uconv.ConvID,
// 		uconv.UserID,
// 	}

// 	conv, err := testQueries.GetUser_conversation(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, conv)
// 	require.Equal(t, arg.ConvID, conv.ConvID)
// 	require.Equal(t, arg.UserID, conv.UserID)
// }

// func TestListUserConvsByUser(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	user2 := createRandomUser(t)
// 	for i := 0; i < 5; i++ {
// 		conv := createRandConv(t, user1.ID, user2.ID)
// 		createRndUsrConv(t, user1.ID, conv.ID)
// 	}

// 	list, err := testQueries.ListUser_conversationByUser(context.Background(), user1.ID)

// 	require.NoError(t, err)
// 	require.Len(t, list, 5)

// 	for _, items := range list {
// 		require.NotEmpty(t, items)
// 		require.Equal(t, user1.ID, items.UserID)
// 	}
// }

