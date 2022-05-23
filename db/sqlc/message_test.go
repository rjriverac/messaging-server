package db

// import (
// 	"context"
// 	"database/sql"
// 	"testing"
// 	"time"

// 	"github.com/rjriverac/messaging-server/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandMessage(t *testing.T, uid int64) Message {

// 	arg := CreateMessageParams{
// 		UserID:  uid,
// 		Content: util.NullStrGen(80),
// 	}
// 	message, err := testQueries.CreateMessage(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, message)

// 	require.Equal(t, arg.UserID, message.UserID)
// 	require.Equal(t, arg.Content.String, message.Content.String)
// 	require.NotEmpty(t, message.CreatedAt)

// 	return message
// }

// func TestCreateMessage(t *testing.T) {
// 	user := createRandomUser(t)
// 	createRandMessage(t, user.ID)
// }

// func TestGetMessage(t *testing.T) {
// 	user := createRandomUser(t)

// 	msg1 := createRandMessage(t, user.ID)
// 	msg2, err := testQueries.GetMessage(context.Background(), msg1.ID)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, msg2)

// 	require.Equal(t, msg1.Content, msg2.Content)
// 	require.Equal(t, msg1.ID, msg2.ID)
// 	require.Equal(t, msg1.UserID, msg2.UserID)
// 	require.WithinDuration(t, msg1.CreatedAt, msg2.CreatedAt, time.Second)
// }

// func TestListMessages(t *testing.T) {
// 	user := createRandomUser(t)

// 	for i := 0; i < 20; i++ {
// 		createRandMessage(t, user.ID)
// 	}

// 	messages, err := testQueries.ListMessageByUser(context.Background(), user.ID)

// 	require.NoError(t, err)
// 	require.Len(t, messages, 20)

// 	for _, msg := range messages {
// 		require.NotEmpty(t, msg)
// 	}

// }

// func TestDeleteMessage(t *testing.T) {
// 	user := createRandomUser(t)

// 	msg := createRandMessage(t, user.ID)
// 	err := testQueries.DeleteMessage(context.Background(),msg.ID)
// 	require.NoError(t,err)

// 	msg2, err := testQueries.GetMessage(context.Background(), msg.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, msg2)
// }
