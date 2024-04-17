package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"learning.com/golang_backend/utils"
)

func TestPasetoToken(t *testing.T) {
	token, err := NewPasetoToken(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	tokenString, err := token.Create(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := token.Verify(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	token, err := NewPasetoToken(utils.RandomString(32))
	require.NoError(t, err)

	tokenString, err := token.Create(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := token.Verify(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
