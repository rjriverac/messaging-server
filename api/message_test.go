package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/rjriverac/messaging-server/db/mock"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {

	msgParams := randomMsgParams()
	result := randomSendResult()

	expectedArgs := db.SendMessageParams{
		CreateMessageParams: db.CreateMessageParams{
			From: msgParams.From,
			Content: msgParams.Content,
			ConvID: msgParams.ConvID,
		},
		UserID: msgParams.UserID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
	SendMessage(gomock.Any(),expectedArgs).
	Times(1).
	Return(result,nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := "/message"

	marshalled, _ := json.Marshal(msgParams)

	request, err := http.NewRequest(http.MethodPost,url,bytes.NewReader(marshalled))
	require.NoError(t,err)

	server.router.ServeHTTP(recorder,request)
	require.Equal(t,http.StatusAccepted,recorder.Code)

}

func randomMsgParams() NewMessageReq {
	return NewMessageReq{
		From: util.RandomString(10),
		Content: util.RandomString(10),
		ConvID: util.RandomInt(0,1000),
		UserID: util.RandomInt(0,1000),
	}
}

func randomSendResult() db.SendResult {
	now := time.Now()
	return db.SendResult{
		Timestamp: now,
		MsgID: util.RandomInt(0,1000),
	}
}