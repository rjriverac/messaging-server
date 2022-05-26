package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func TestGetUser(t *testing.T) {
	account := randomUser()

	testCases := []struct {
		name       string
		uID        int64
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			uID:  account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserBody(t, recorder.Body, account)
			},
		},
		{
			name: "NotFound",
			uID:  account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.GetUserRow{}, sql.ErrNoRows)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IntServerErr",
			uID:  account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.GetUserRow{}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			uID:  0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", tc.uID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t,recorder)
		})
	}
}

func requireUserBody(t *testing.T, res *bytes.Buffer, user db.GetUserRow) {
	data, err := ioutil.ReadAll(res)
	require.NoError(t, err)

	var responseAcc db.GetUserRow
	err = json.Unmarshal(data, &responseAcc)
	require.NoError(t, err)
	require.Equal(t, responseAcc.Email, user.Email)
	require.Equal(t, responseAcc.ID, user.ID)
	require.Equal(t, responseAcc.Image, user.Image)
	require.Equal(t, responseAcc.Name, user.Name)
	require.Equal(t, responseAcc.Status, user.Status)
	require.WithinDuration(t, responseAcc.CreatedAt, user.CreatedAt, time.Second)
}

func randomUser() db.GetUserRow {
	now := time.Now()
	time := now.Add(time.Duration(-10) * time.Minute)
	return db.GetUserRow{
		ID:        util.RandomInt(1, 1000),
		Name:      util.RandomUserGen(),
		Email:     util.RandomEmail(),
		Image:     sql.NullString{String: util.RandomString(10), Valid: true},
		Status:    sql.NullString{String: util.RandomString(10), Valid: true},
		CreatedAt: time,
	}
}
