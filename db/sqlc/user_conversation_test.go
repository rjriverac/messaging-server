package db

import (
	"context"
	"database/sql"

	"testing"

	"github.com/stretchr/testify/require"
)

func createRndUsrConv(t *testing.T) UserConversation {

	user := createRandomUser(t)
	conv := createRandConv(t)

	arg := CreateUser_conversationParams{
		UserID: user.ID,
		ConvID: conv.ID,
	}

	uconv, err := testQueries.CreateUser_conversation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, uconv)
	require.Equal(t, arg.UserID, uconv.UserID)
	require.Equal(t, arg.UserID, uconv.UserID)

	return uconv
}

func TestCreateUsrConv(t *testing.T) {
	createRndUsrConv(t)
}

func TestGetUserConv(t *testing.T) {
	nconv := createRndUsrConv(t)

	arg := GetUser_conversationParams{
		UserID: nconv.UserID,
		ConvID: nconv.ConvID,
	}
	conv, err := testQueries.GetUser_conversation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, conv)
	require.Equal(t, arg.ConvID, conv.ConvID)
	require.Equal(t, arg.UserID, conv.UserID)
}

func TestListUserConvs(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRndUsrConv(t)
	}
	lconv, err := testQueries.ListUser_conversations(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, lconv)
	for _, v := range lconv {
		require.NotEmpty(t, v)
	}
}

func TestListUConvByUser(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 10; i++ {
		conv := createRandConv(t)

		arg := CreateUser_conversationParams{
			UserID: user.ID,
			ConvID: conv.ID,
		}
		testQueries.CreateUser_conversation(context.Background(), arg)
	}

	lconv,err := testQueries.ListUser_conversationByUser(context.Background(),user.ID)
	require.NoError(t,err)
	require.NotEmpty(t,lconv)
	require.Len(t,lconv,10)
	for _, c := range lconv {
		require.NotEmpty(t,c)
	}
}

func TestDeleteUsrConv(t *testing.T) {
	conv := createRndUsrConv(t)

	arg := DeleteUser_conversationParams{
		UserID: conv.UserID,
		ConvID: conv.ConvID,
	}
	err := testQueries.DeleteUser_conversation(context.Background(),arg)
	require.NoError(t,err)

	narg := GetUser_conversationParams{UserID: conv.UserID,ConvID: conv.ConvID}
	conv2, lerr := testQueries.GetUser_conversation(context.Background(),narg)
	require.Error(t,lerr)
	require.Empty(t,conv2)
	require.EqualError(t,lerr,sql.ErrNoRows.Error())
}

func TestGetUConvById(t *testing.T) {
	conv := createRndUsrConv(t)

	conv2,err := testQueries.GetUser_conv_by_id(context.Background(),conv.ID)
	require.NoError(t,err)
	require.NotEmpty(t,conv2)
	require.Equal(t,conv.ConvID,conv2.ConvID)
	require.Equal(t,conv.UserID,conv2.UserID)
	require.Equal(t,conv.ID,conv2.ID)
}

func TestDelUConvById(t *testing.T) {
	conv := createRndUsrConv(t)

	err := testQueries.DeleteUser_conversation_by_id(context.Background(),conv.ID)
	require.NoError(t,err)
	conv2,nerr :=testQueries.GetUser_conv_by_id(context.Background(),conv.ID)
	require.Empty(t,conv2)
	require.Error(t,nerr)
	require.EqualError(t,nerr,sql.ErrNoRows.Error())
}