package api

import (
	"bytes"
	"database/sql"
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

	testCases := []struct {
		name       string
		arg        db.SendMessageParams
		buildStubs func(store *mockdb.MockStore, params db.SendMessageParams)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		arg: db.SendMessageParams{
			CreateMessageParams: db.CreateMessageParams{
				From:    msgParams.From,
				Content: msgParams.Content,
				ConvID:  msgParams.ConvID,
			},
			UserID: msgParams.UserID,
		},
		buildStubs: func(store *mockdb.MockStore, params db.SendMessageParams) {
			store.EXPECT().
				SendMessage(gomock.Any(), gomock.Eq(params)).
				Times(1).
				Return(result, nil)
		},
		checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusAccepted, recorder.Code)
		},
	}, {
		name: "Bad Request",
		arg: db.SendMessageParams{
			CreateMessageParams: db.CreateMessageParams{
				From:    "",
				Content: msgParams.Content,
				ConvID:  msgParams.ConvID,
			},
			UserID: 0,
		},
		buildStubs: func(store *mockdb.MockStore, params db.SendMessageParams) {
			store.EXPECT().
				SendMessage(gomock.Any(), params).
				Times(0)
		},
		checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	}, {
		name: "Int Server Err",
		arg: db.SendMessageParams{
			CreateMessageParams: db.CreateMessageParams{
				From:    msgParams.From,
				Content: msgParams.Content,
				ConvID:  msgParams.ConvID,
			},
			UserID: msgParams.UserID,
		},
		buildStubs: func(store *mockdb.MockStore, params db.SendMessageParams) {
			store.EXPECT().
				SendMessage(gomock.Any(), gomock.Any()).
				Times(1).
				Return(result, sql.ErrConnDone)
		},
		checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	}}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store, tc.arg)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/message"
			marshalled, _ := json.Marshal(tc.arg)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkRes(t, recorder)
		})
	}

}

func randomMsgParams() NewMessageReq {
	return NewMessageReq{
		From:    util.RandomString(10),
		Content: util.RandomString(10),
		ConvID:  util.RandomInt(0, 1000),
		UserID:  util.RandomInt(0, 1000),
	}
}

func randomSendResult() db.SendResult {
	now := time.Now()
	return db.SendResult{
		Timestamp: now,
		MsgID:     util.RandomInt(0, 1000),
	}
}
