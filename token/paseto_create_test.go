package token

import (
	"testing"
	"time"

	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	id := util.RandomInt(1, 1000)
	duration := time.Minute

	issuedAt := time.Now()
	expires := time.Now().Add(duration)

	token, err := maker.CreateToken(id, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, id, payload.User)
	require.WithinDuration(t, issuedAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expires, payload.Expires, time.Second)
}
func TestPasetoExpired(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomInt(0, 1000), -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
