package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) CreateUserRow {
	arg := CreateUserParams{
		Name:     util.RandomUserGen(),
		Email:    util.RandomEmail(),
		HashedPw: util.RandomHashedPW(),
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

//! fix PW being passed in as empty string in update func
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
	fmt.Printf("user3: %v status: %v",user3,partialarg.Status)
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
