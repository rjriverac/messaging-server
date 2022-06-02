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
			tc.checkRes(t, recorder)
		})
	}
}

func TestListUsers(t *testing.T) {

	list := randomListUser(10)

	testCases := []struct {
		name       string
		params     db.ListUsersParams
		buildStubs func(store *mockdb.MockStore, params db.ListUsersParams)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: db.ListUsersParams{Limit: 10, Offset: 1},
			buildStubs: func(store *mockdb.MockStore, params db.ListUsersParams) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(db.ListUsersParams{Limit: params.Limit, Offset: (params.Offset - 1) * params.Limit})).
					Times(1).
					Return(list, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				listMatch(t, recorder.Body, list)
			},
		}, {
			name:   "Bad Request",
			params: db.ListUsersParams{Limit: 50, Offset: 1},
			buildStubs: func(store *mockdb.MockStore, params db.ListUsersParams) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:   "Internal Server Err",
			params: db.ListUsersParams{Limit: 10, Offset: 1},
			buildStubs: func(store *mockdb.MockStore, params db.ListUsersParams) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(db.ListUsersParams{Limit: params.Limit, Offset: (params.Offset - 1) * params.Limit})).
					Times(1).
					Return(list, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store, tc.params)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account?page_id=%d&page_size=%d", tc.params.Offset, tc.params.Limit)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder)
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

func randomListUser(num int) []db.ListUsersRow {
	var list []db.ListUsersRow
	for i := 0; i < num; i++ {
		list = append(list, db.ListUsersRow{
			ID:     util.RandomInt(1, 1000),
			Name:   util.RandomUserGen(),
			Email:  util.RandomEmail(),
			Image:  util.NullStrGen(10),
			Status: util.NullStrGen(15),
		})
	}
	return list
}

func listMatch(t *testing.T, res *bytes.Buffer, list []db.ListUsersRow) {
	data, err := ioutil.ReadAll(res)
	require.NoError(t, err)

	var retrievedList []db.ListUsersRow
	err = json.Unmarshal(data, &retrievedList)
	require.NoError(t, err)
	require.Equal(t, len(list), len(retrievedList))
	for _, row := range list {
		for _, qrow := range retrievedList {
			if row == qrow {
				require.Equal(t, row.Email, qrow.Email)
				require.Equal(t, row.ID, qrow.ID)
				require.Equal(t, row.Image, qrow.Image)
				require.Equal(t, row.Name, qrow.Name)
				require.Equal(t, row.Status, qrow.Status)
			}
		}
	}
}
func TestCreateUser(t *testing.T) {

	cUserParams := db.CreateUserParams{
		Name:     util.RandomUserGen(),
		Email:    util.RandomEmail(),
		HashedPw: util.RandomHashedPW(),
	}

	anotherUser := db.CreateUserRow{
		ID:        util.RandomInt(0, 1000),
		Name:      cUserParams.Name,
		Email:     cUserParams.Email,
		Image:     sql.NullString{Valid: false},
		Status:    sql.NullString{Valid: false},
		CreatedAt: time.Now(),
	}

	testCases := []struct {
		name       string
		params     db.CreateUserParams
		buildStubs func(store *mockdb.MockStore, params db.CreateUserParams)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: db.CreateUserParams{
				Name:     util.RandomUserGen(),
				Email:    util.RandomEmail(),
				HashedPw: util.RandomHashedPW(),
			},
			buildStubs: func(store *mockdb.MockStore, params db.CreateUserParams) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(anotherUser, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "bad request",
			params: db.CreateUserParams{
				Name: "",
				Email: "",
				HashedPw: "",
			},
			buildStubs: func(store *mockdb.MockStore,params db.CreateUserParams) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
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
			tc.buildStubs(store,tc.params)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/account"

			marshalled, _ := json.Marshal(tc.params)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder)
		})
	}
}
