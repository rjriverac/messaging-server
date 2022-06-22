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
	"github.com/lib/pq"
	mockdb "github.com/rjriverac/messaging-server/db/mock"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
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

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"name":     user.Name,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.CreateUserParams{
					Name:  user.Name,
					Email: user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
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
					Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "bad request",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": 3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "duplicate email",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		}, {
			name: "short password",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": "123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "invalid email",
			body: gin.H{
				"name":     user.Name,
				"email":    "email-com",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/account/"

			marshalled, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			// requestDump, err := httputil.DumpRequest(request, true)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Printf("\n****request dump****\n,%v\n\n", string(requestDump))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(recorder)
		})
	}
}

func TestGetUser(t *testing.T) {
	account := randomUser()
	mismatch := account.ID + 1

	testCases := []struct {
		name       string
		uID        int64
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, account.ID, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserBody(t, recorder.Body, account)
			},
		}, {
			name: "Mismatched IDs",
			uID:  account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, mismatch, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, account.ID, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, account.ID, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, account.ID, time.Minute)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", tc.uID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

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
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore, params db.ListUsersParams)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: db.ListUsersParams{Limit: 10, Offset: 1},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, list[0].ID, time.Minute)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, list[0].ID, time.Minute)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, list[0].ID, time.Minute)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/?page_id=%d&page_size=%d", tc.params.Offset, tc.params.Limit)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

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

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserInfoParams
	password string
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserInfoParams)

	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPw.String)
	if err != nil {
		return false
	}
	e.arg.HashedPw = arg.HashedPw
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("Matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateParamsUser(arg db.UpdateUserInfoParams, password string) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password}
}

func TestUpdateUser(t *testing.T) {

	user, _ := randomDBUser(t)
	now := time.Now()

	name := ToBeNullString(util.RandomUserGen())
	email := ToBeNullString(util.RandomEmail())
	image := ToBeNullString(util.RandomString(10))
	status := ToBeNullString(util.RandomString(10))
	newpw := ToBeNullString(util.RandomString(8))

	testCases := []struct {
		name       string
		uId        int64
		body       gin.H
		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore, uId int64)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64)
	}{
		{
			name: "OK",
			uId:  user.ID,
			body: gin.H{
				"name":     name,
				"email":    email,
				"image":    image,
				"status":   status,
				"password": newpw,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, uId int64) {
				arg := db.UpdateUserInfoParams{
					Name:   name.ToNstring(),
					Email:  email.ToNstring(),
					Image:  image.ToNstring(),
					Status: status.ToNstring(),
					ID:     uId,
				}

				store.EXPECT().
					UpdateUserInfo(gomock.Any(), EqUpdateParamsUser(arg, string(newpw))).
					Times(1).
					Return(db.UpdateUserInfoRow{
						ID:        uId,
						Name:      string(name),
						Email:     string(email),
						Image:     image.ToNstring(),
						Status:    status.ToNstring(),
						CreatedAt: now,
					}, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireUserUpdateBody(t, recorder.Body, req, uID)
			},
		}, {
			name: "No Pw Ok",
			uId:  user.ID,
			body: gin.H{
				"name":   name,
				"email":  email,
				"image":  image,
				"status": status,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, uId int64) {
				arg := db.UpdateUserInfoParams{
					Name:   name.ToNstring(),
					Email:  email.ToNstring(),
					Image:  image.ToNstring(),
					Status: status.ToNstring(),
					ID:     uId,
				}

				store.EXPECT().
					UpdateUserInfo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.UpdateUserInfoRow{
						ID:        uId,
						Name:      string(name),
						Email:     string(email),
						Image:     image.ToNstring(),
						Status:    status.ToNstring(),
						CreatedAt: now,
					}, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
				requireUserUpdateBody(t, recorder.Body, req, uID)
			},
		}, {
			name: "token error",
			uId:  0,
			body: gin.H{
				"name":     name,
				"email":    email,
				"image":    image,
				"status":   status,
				"password": newpw,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, "token", 0, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, uId int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		}, {
			name: "bad request json",
			uId:  user.ID,
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, uId int64) {
				store.EXPECT().
					UpdateUserInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "Internal Server Err",
			uId:  user.ID,
			body: gin.H{
				"name":     name,
				"email":    email,
				"image":    image,
				"status":   status,
				"password": newpw,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, uId int64) {
				arg := db.UpdateUserInfoParams{
					Name:   name.ToNstring(),
					Email:  email.ToNstring(),
					Image:  image.ToNstring(),
					Status: status.ToNstring(),
					ID:     uId,
				}

				store.EXPECT().
					UpdateUserInfo(gomock.Any(), EqUpdateParamsUser(arg, string(newpw))).
					Times(1).
					Return(db.UpdateUserInfoRow{
						ID:        uId,
						Name:      string(name),
						Email:     string(email),
						Image:     image.ToNstring(),
						Status:    status.ToNstring(),
						CreatedAt: now,
					}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder, req gin.H, uID int64) {
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
			tc.buildStubs(store, tc.uId)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/account/"

			marshalled, _ := json.Marshal(tc.body)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder, tc.body, tc.uId)

		})
	}

}

func requireUserUpdateBody(t *testing.T, res *bytes.Buffer, req gin.H, uID int64) {
	data, err := ioutil.ReadAll(res)
	require.NoError(t, err)
	now := time.Now()

	var user UpdateUserReturn
	err = json.Unmarshal(data, &user)

	require.NoError(t, err)
	require.Equal(t, string(req["email"].(ToBeNullString)), user.Email)
	require.Equal(t, string(req["name"].(ToBeNullString)), user.Name)
	require.Equal(t, string(req["status"].(ToBeNullString)), user.Status)
	require.Equal(t, string(req["image"].(ToBeNullString)), user.Image)
	require.Equal(t, uID, user.ID)
	require.WithinDuration(t, now, user.CreatedAt, time.Second)
}

func TestLoginUser(t *testing.T) {

	user, password := randomDBUser(t)
	testCases := []struct {
		desc       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			desc: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			desc: "Bad Request",
			body: gin.H{
				"email":    5,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			desc: "User not found",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			desc: "Internal Server Err",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			desc: "Wrong PW",
			body: gin.H{
				"email":    user.Email,
				"password": "badpassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tC.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/account/login"

			marshalled, _ := json.Marshal(tC.body)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tC.checkRes(t, recorder)
		})
	}
}
