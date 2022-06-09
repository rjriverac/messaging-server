package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashedPw1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPw1)

	err = CheckPassword(password, hashedPw1)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	err = CheckPassword(wrongPassword, hashedPw1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPw2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPw2)
	require.NotEqual(t, hashedPw1, hashedPw2)
}
