package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"learning.com/golang_backend/utils"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	// trouver une moyen de faire mieux par ici
	account, err := testRepository.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestDeleteAccount(t *testing.T) {
	expected := createRandomAccount(t)
	err := testRepository.DeleteAccount(context.Background(), expected.ID)
	require.NoError(t, err)

	result, err := testRepository.GetAccount(context.Background(), expected.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, result)
}

func TestGetAccount(t *testing.T) {
	expected := createRandomAccount(t)
	result, err := testRepository.GetAccount(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Owner, result.Owner)
	require.Equal(t, expected.Currency, result.Currency)
	require.Equal(t, expected.Balance, result.Balance)
	require.WithinDuration(t, expected.CreatedAt.Time, result.CreatedAt.Time, time.Second)
}

func TestListAcounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAcountsParams{
		Offset: 0,
		Limit:  5,
	}

	accounts, err := testRepository.ListAcounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	expected := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      expected.ID,
		Balance: utils.RandomMoney(),
	}
	result, err := testRepository.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Owner, result.Owner)
	require.Equal(t, expected.Currency, result.Currency)
	require.Equal(t, arg.Balance, result.Balance)
	require.WithinDuration(t, expected.CreatedAt.Time, result.CreatedAt.Time, time.Second)
}
