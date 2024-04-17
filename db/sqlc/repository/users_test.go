package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"learning.com/golang_backend/utils"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	ownerFullName := utils.RandomOwner()
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       ownerFullName,
		Email:          utils.OwnerEmail(ownerFullName),
	}

	// trouver une moyen de faire mieux par ici
	user, err := testRepository.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	expected := createRandomUser(t)
	result, err := testRepository.GetUser(context.Background(), expected.Username)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.Username, result.Username)
	require.Equal(t, expected.HashedPassword, result.HashedPassword)
	require.Equal(t, expected.FullName, result.FullName)
	require.Equal(t, expected.Email, result.Email)
	require.WithinDuration(t, expected.PasswordChangedAt, result.PasswordChangedAt, time.Second)
	require.WithinDuration(t, expected.CreatedAt, result.CreatedAt, time.Second)
}
