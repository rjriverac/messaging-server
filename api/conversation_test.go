package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/rjriverac/messaging-server/db/mock"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestGetConvos(t *testing.T) {

	user, _ := randomDBUser(t)

	n := 5
	convs := make([]db.Conversation, n)

	for i := 0; i < n; i++ {
		convs[i] = db.Conversation{
			ID:   util.RandomInt(1, 1000),
			Name: util.NullStrGen(6),
		}
	}
	convs[4].Name = sql.NullString{Valid: false}

	testCases := []struct {
		name       string
		buildStubs func(store *mockdb.MockStore)
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListConvFromUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(convs, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Found",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListConvFromUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(convs, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListConvFromUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(convs, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tC.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/conversation"

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			tC.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tC.checkRes(t, recorder)
		})
	}
}
