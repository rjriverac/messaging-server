package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rjriverac/messaging-server/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)

	id := util.RandomInt(1, 1000)
	duration := time.Minute

	issuedAt := time.Now()
	expires := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(id, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, id, payload.User)
	require.WithinDuration(t, issuedAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expires, payload.Expires, time.Second)
}
func TestJWTExpired(t *testing.T) {
	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomInt(0, 1000), -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTokenBadHeader(t *testing.T) {
	payload, err := CreatePayload(util.RandomInt(0, 1000), time.Hour)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
