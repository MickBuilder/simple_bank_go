package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"learning.com/golang_backend/utils"
)

// Helper function to create a random entry
func createRandomEntry(t *testing.T, accountID int64) Entry {
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testRepository.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	expected := createRandomEntry(t, account.ID)

	result, err := testRepository.GetEntry(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.AccountID, result.AccountID)
	require.Equal(t, expected.Amount, result.Amount)
	require.WithinDuration(t, expected.CreatedAt, result.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account.ID)
	}

	arg := ListEntriesParams{
		Offset: 0,
		Limit:  5,
	}

	entries, err := testRepository.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
