package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/rjriverac/messaging-server/db/mock"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {

	msgParams := randomMsgParams()
	result := randomSendResult()
	user, _ := randomDBUser(t)

	testCases := []struct {
		name       string
		arg        gin.H
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
		},
		arg: gin.H{
			"content": msgParams.Content,
			"convID":  msgParams.ConvID,
		},
		buildStubs: func(store *mockdb.MockStore) {
			// store.EXPECT().
			// 	SendMessage(gomock.Any(), gomock.Any()).
			// 	Times(1).
			// 	Return(result, nil)
			gomock.InOrder(
				store.EXPECT().SendMessage(gomock.Any(), db.SendMessageParams{
					CreateMessageParams: db.CreateMessageParams{
						From:    user.Name,
						Content: msgParams.Content,
						ConvID:  msgParams.ConvID,
					},
					UserID: user.ID,
				}).Times(1).Return(result, nil),
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1),
			)

		},
		checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusAccepted, recorder.Code)
		},
	}, {
		name: "Bad Request",
		arg: gin.H{
			"content": 0,
			"convID":  msgParams.ConvID,
		},
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
		},
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				SendMessage(gomock.Any(), gomock.Any()).
				Times(0)
		},
		checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	}, {
		name: "Int Server Err",
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
		},
		arg: gin.H{},
		buildStubs: func(store *mockdb.MockStore) {
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
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/message"
			marshalled, _ := json.Marshal(msgParams)

			// fmt.Printf("marshalled request:%c", marshalled)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))

			dump, _ := httputil.DumpRequest(request, false)

			fmt.Printf("dumped request:%v\n", string(dump))
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkRes(t, recorder)
		})
	}

}

func randomMsgParams() NewMessageReq {
	return NewMessageReq{
		// From:    util.RandomString(10),
		Content: util.RandomString(10),
		ConvID:  util.RandomInt(0, 1000),
		// UserID:  util.RandomInt(0, 1000),
	}
}

func randomSendResult() db.SendResult {
	now := time.Now()
	return db.SendResult{
		Timestamp: now,
		MsgID:     util.RandomInt(0, 1000),
	}
}
