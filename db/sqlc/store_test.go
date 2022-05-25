package db

import (
	"context"
	"testing"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	store := NewStore(testDB)

	sender := createRandomUser(t)
	message := CreateMessageParams{
		From:    sender.Name,
		Content: util.RandomString(50),
		ConvID:  createRandConv(t).ID,
	}

	n := 5

	errs := make(chan error)
	res := make(chan SendResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.SendMessage(
				context.Background(),
				SendMessageParams{UserID: sender.ID, CreateMessageParams: &message},
			)
			errs <- err
			res <- result

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-res
		require.NotEmpty(t, result)
		require.NotZero(t, result.Timestamp)

		_, err = store.GetMessage(context.Background(), result.MsgID)
		require.NoError(t, err)

		_, err = store.GetUser_conversation(context.Background(), GetUser_conversationParams{UserID: sender.ID, ConvID: message.ConvID})
		require.NoError(t,err)
	}
}
