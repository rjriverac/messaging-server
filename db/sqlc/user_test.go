package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPw, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Name:     util.RandomUserGen(),
		Email:    util.RandomEmail(),
		HashedPw: hashedPw,
		Image:    util.NullStrGen(8),
		Status:   sql.NullString{String: "", Valid: false},
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Image, user.Image)
	require.Equal(t, arg.Status, user.Status)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Image, user2.Image)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Status, user2.Status)
	require.Equal(t, user.ID, user2.ID)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
}
func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Image, user2.Image)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Status, user2.Status)
	require.Equal(t, user.HashedPw, user2.HashedPw)
	require.Equal(t, user.ID, user2.ID)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserInfoParams{
		ID:       user.ID,
		Name:     util.NullStrGen(5),
		Email:    sql.NullString{String: util.RandomEmail(), Valid: true},
		Image:    util.NullStrGen(15),
		Status:   util.NullStrGen(15),
		HashedPw: sql.NullString{String: "", Valid: false},
	}

	user2, err := testQueries.UpdateUserInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, arg.ID, user2.ID)
	require.Equal(t, arg.Email.String, user2.Email)
	require.Equal(t, arg.Image.String, user2.Image.String)
	require.Equal(t, arg.Name.String, user2.Name)
	require.Equal(t, arg.Status.String, user2.Status.String)

	partialarg := UpdateUserInfoParams{
		ID:     user.ID,
		Status: util.NullStrGen(5),
	}

	user3, err := testQueries.UpdateUserInfo(context.Background(), partialarg)
	require.NoError(t, err)
	require.NotEmpty(t, user3)
	require.Equal(t, arg.ID, user3.ID)
	require.Equal(t, arg.Email.String, user3.Email)
	require.Equal(t, arg.Image.String, user3.Image.String)
	require.Equal(t, arg.Name.String, user3.Name)
	require.Equal(t, partialarg.Status.String, user3.Status.String)
	require.Len(t, user3.Status.String, 5)

}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestListUserMessages(t *testing.T) {

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	convs := make([]Conversation, 2)

	conv1 := createRandConv(t)
	conv2 := createRandConv(t)
	convs = append(convs, conv1, conv2)
	for _, conv := range convs {
		for i := 0; i < 2; i++ {
			arg := CreateMessageParams{
				From:    user1.Name,
				Content: util.RandomString(20),
				ConvID:  conv.ID,
			}
			arg2 := CreateMessageParams{
				From:    user2.Name,
				Content: util.RandomString(20),
				ConvID:  conv.ID,
			}
			testQueries.CreateMessage(context.Background(), arg)
			testQueries.CreateMessage(context.Background(), arg2)
		}
		testQueries.CreateUser_conversation(
			context.Background(),
			CreateUser_conversationParams{
				UserID: user1.ID,
				ConvID: conv.ID,
			},
		)
		testQueries.CreateUser_conversation(
			context.Background(),
			CreateUser_conversationParams{
				UserID: user2.ID,
				ConvID: conv.ID,
			},
		)
	}

	u1Msg, err1 := testQueries.ListUserMessages(context.Background(), user1.ID)
	u2Msg, err2 := testQueries.ListUserMessages(context.Background(), user2.ID)

	require.NoError(t, err1)
	require.NoError(t, err2)

	require.NotEmpty(t, u1Msg)
	require.NotEmpty(t, u2Msg)

	require.Len(t, u1Msg, 8)
	require.Len(t, u2Msg, 8)
	for _, msg := range u1Msg {
		require.NotEmpty(t, msg)
	}
	for _, msg := range u2Msg {
		require.NotEmpty(t, msg)
	}

}
func TestListConvFromUser(t *testing.T) {
	user1 := createRandomUser(t)
	var convs []Conversation
	for i := 0; i < 10; i++ {
		c := createRandConv(t)
		convs = append(convs, c)
	}
	for _, con := range convs {
		usr := createRandomUser(t)
		for i := 0; i < 5; i++ {
			if i%2 != 0 {
				arg := CreateMessageParams{
					From:    usr.Name,
					Content: util.RandomString(20),
					ConvID:  con.ID,
				}
				testQueries.CreateMessage(context.Background(), arg)
			} else {
				arg := CreateMessageParams{
					From:    user1.Name,
					Content: util.RandomString(20),
					ConvID:  con.ID,
				}
				testQueries.CreateMessage(context.Background(), arg)
			}
		}
		testQueries.CreateUser_conversation(context.Background(), CreateUser_conversationParams{
			UserID: user1.ID,
			ConvID: con.ID,
		})
		testQueries.CreateUser_conversation(context.Background(), CreateUser_conversationParams{
			UserID: usr.ID,
			ConvID: con.ID,
		})
	}

	list, err := testQueries.ListConvFromUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.Len(t, list, 10)
	for _, c := range list {
		require.NotEmpty(t, c)
		require.NotEmpty(t, c.Name)
	}

}
