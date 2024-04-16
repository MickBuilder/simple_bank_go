package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"learning.com/golang_backend/utils"
)

// Helper function to create a random transfer
func createRandomTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testRepository.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	expected := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

	result, err := testRepository.GetTransfer(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.FromAccountID, result.FromAccountID)
	require.Equal(t, expected.ToAccountID, result.ToAccountID)
	require.Equal(t, expected.Amount, result.Amount)
	require.WithinDuration(t, expected.CreatedAt, result.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
		createRandomTransfer(t, toAccount.ID, fromAccount.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Offset:        0,
		Limit:         5,
	}

	transfers, err := testRepository.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
