package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/rjriverac/messaging-server/db/mock"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPw)
	if err != nil {
		return false
	}
	e.arg.HashedPw = arg.HashedPw
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("Matches arg %v and password %v", e.arg, e.password)
}
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {

	user, password := randomDBUser(t)
	user.Email = fmt.Sprint(util.RandomString(5), "@email.com")
	fmt.Printf("user: %v \n password: %v\n hashedPw: %v\n", user, password, user.HashedPw)

	anotherUser := db.CreateUserRow{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Image:     sql.NullString{Valid: false},
		Status:    sql.NullString{Valid: false},
		CreatedAt: user.CreatedAt,
	}

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Name:   user.Name,
					Email:  user.Email,
					Image:  user.Image,
					Status: user.Status,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(anotherUser, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name: "int server err",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.CreateUserRow{}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "bad request",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
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
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/account"

			marshalled, err := json.Marshal(tc.body)
			require.NoError(t, err)

			fmt.Printf("marshalled:%v\n", string(marshalled))

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))

			buf := new(bytes.Buffer)
			buf.ReadFrom(request.Body)
			newStr := buf.String()

			fmt.Printf("request:%v", newStr)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder)
		})
	}
}

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

	var responseAcc GetUserReturn
	err = json.Unmarshal(data, &responseAcc)
	require.NoError(t, err)
	require.Equal(t, responseAcc.Email, user.Email)
	require.Equal(t, responseAcc.ID, user.ID)
	require.Equal(t, responseAcc.Image, user.Image.String)
	require.Equal(t, responseAcc.Name, user.Name)
	require.Equal(t, responseAcc.Status, user.Status.String)
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

func randomDBUser(t *testing.T) (user db.User, password string) {
	password = util.RandomHashedPW()
	hashed, err := util.HashPassword(password)
	require.NoError(t, err)

	now := time.Now()
	time := now.Add(time.Duration(-10) * time.Minute)

	user = db.User{
		ID:        util.RandomInt(1, 1000),
		Name:      util.RandomUserGen(),
		Email:     util.RandomEmail(),
		HashedPw:  hashed,
		Image:     sql.NullString{String: util.RandomString(10), Valid: true},
		Status:    sql.NullString{String: util.RandomString(10), Valid: true},
		CreatedAt: time,
	}
	return
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

	var retrievedList []ListUserAcc
	err = json.Unmarshal(data, &retrievedList)
	require.NoError(t, err)
	require.Equal(t, len(list), len(retrievedList))
	for _, row := range list {
		for _, qrow := range retrievedList {
			if row.ID == qrow.ID {
				require.Equal(t, row.Email, qrow.Email)
				require.Equal(t, row.ID, qrow.ID)
				require.Equal(t, row.Image.String, qrow.Image)
				require.Equal(t, row.Name, qrow.Name)
				require.Equal(t, row.Status.String, qrow.Status)
			}
		}
	}
}

func TestUpdateUser(t *testing.T) {

	user := randomUser()
	now := time.Now()

	testCases := []struct {
		name       string
		uId        int64
		params     UpdateUserRequest
		buildStubs func(store *mockdb.MockStore, params UpdateUserRequest, uID int64)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder, params UpdateUserRequest, uID int64)
	}{
		{
			name: "OK",
			uId:  user.ID,
			params: UpdateUserRequest{
				Name:     ToBeNullString(util.RandomUserGen()),
				Email:    ToBeNullString(util.RandomEmail()),
				Image:    ToBeNullString(util.RandomString(5)),
				Status:   ToBeNullString(util.RandomString(10)),
				HashedPw: ToBeNullString(""),
			},
			buildStubs: func(store *mockdb.MockStore, params UpdateUserRequest, uId int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), db.UpdateUserInfoParams{
						Name:     params.Name.Scan(params.Name),
						Email:    params.Email.Scan(params.Email),
						Image:    params.Image.Scan(params.Image),
						Status:   params.Status.Scan(params.Status),
						HashedPw: params.HashedPw.Scan(params.HashedPw),
						ID:       uId,
					}).
					Times(1).
					Return(db.UpdateUserInfoRow{
						ID:        uId,
						Name:      string(params.Name),
						Email:     string(params.Email),
						Image:     params.Image.Scan(params.Image),
						Status:    params.Status.Scan(params.Status),
						CreatedAt: now,
					}, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, params UpdateUserRequest, uID int64) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireUserUpdateBody(t, recorder.Body, params, uID)
			},
		}, {
			name: "bad request UID",
			uId:  0,
			params: UpdateUserRequest{
				Name:     ToBeNullString(util.RandomUserGen()),
				Email:    ToBeNullString(util.RandomEmail()),
				Image:    ToBeNullString(util.RandomString(5)),
				Status:   ToBeNullString(util.RandomString(10)),
				HashedPw: ToBeNullString(""),
			},
			buildStubs: func(store *mockdb.MockStore, params UpdateUserRequest, uID int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, params UpdateUserRequest, uID int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:   "bad request json",
			uId:    user.ID,
			params: UpdateUserRequest{},
			buildStubs: func(store *mockdb.MockStore, params UpdateUserRequest, uID int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, params UpdateUserRequest, uID int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "Internal Server Err",
			uId:  user.ID,
			params: UpdateUserRequest{
				Name:     ToBeNullString(util.RandomUserGen()),
				Email:    ToBeNullString(util.RandomEmail()),
				Image:    ToBeNullString(util.RandomString(5)),
				Status:   ToBeNullString(util.RandomString(10)),
				HashedPw: ToBeNullString(""),
			},
			buildStubs: func(store *mockdb.MockStore, params UpdateUserRequest, uId int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), db.UpdateUserInfoParams{
						Name:     params.Name.Scan(params.Name),
						Email:    params.Email.Scan(params.Email),
						Image:    params.Image.Scan(params.Image),
						Status:   params.Status.Scan(params.Status),
						HashedPw: params.HashedPw.Scan(params.HashedPw),
						ID:       uId,
					}).
					Times(1).
					Return(db.UpdateUserInfoRow{
						ID:        uId,
						Name:      string(params.Name),
						Email:     string(params.Email),
						Image:     params.Image.Scan(params.Image),
						Status:    params.Status.Scan(params.Status),
						CreatedAt: now,
					}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, params UpdateUserRequest, uID int64) {
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
			tc.buildStubs(store, tc.params, tc.uId)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/?uid=%v", tc.uId)

			marshalled, _ := json.Marshal(tc.params)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder, tc.params, tc.uId)

		})
	}

}

func requireUserUpdateBody(t *testing.T, res *bytes.Buffer, params UpdateUserRequest, uID int64) {
	data, err := ioutil.ReadAll(res)
	require.NoError(t, err)
	now := time.Now()

	var user UpdateUserReturn
	err = json.Unmarshal(data, &user)

	require.NoError(t, err)
	require.Equal(t, string(params.Email), user.Email)
	require.Equal(t, string(params.Name), user.Name)
	require.Equal(t, string(params.Status), user.Status)
	require.Equal(t, string(params.Image), user.Image)
	require.Equal(t, uID, user.ID)
	require.WithinDuration(t, now, user.CreatedAt, time.Second)
}
