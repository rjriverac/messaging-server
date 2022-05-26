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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().GetUser(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireUserBody(t,recorder.Body,account)
}

func requireUserBody(t *testing.T, res *bytes.Buffer, user db.GetUserRow) {
	data, err := ioutil.ReadAll(res)
	require.NoError(t,err)

	var responseAcc db.GetUserRow
	err = json.Unmarshal(data,&responseAcc)
	require.NoError(t,err)
	require.Equal(t,responseAcc.Email,user.Email)
	require.Equal(t,responseAcc.ID,user.ID)
	require.Equal(t,responseAcc.Image,user.Image)
	require.Equal(t,responseAcc.Name,user.Name)
	require.Equal(t,responseAcc.Status,user.Status)
	require.WithinDuration(t,responseAcc.CreatedAt,user.CreatedAt,time.Second)
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
