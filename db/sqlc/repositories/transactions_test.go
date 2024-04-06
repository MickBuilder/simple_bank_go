package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	sourceAccount := createRandomAccount(t)
	destAccount := createRandomAccount(t)
	fmt.Println(">> before:", sourceAccount.Balance, destAccount.Balance)

	tries, amount := 5, int64(10)
	errs, results := make(chan error), make(chan TransferTxResult)

	// run tries concurent transfer transaction
	for i := 0; i < tries; i++ {
		go func() {
			result, err := testRepository.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: sourceAccount.ID,
				ToAccountID:   destAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	// checks results
	for i := 0; i < tries; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, sourceAccount.ID, transfer.FromAccountID)
		require.Equal(t, destAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testRepository.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, sourceAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testRepository.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, destAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testRepository.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, sourceAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, destAccount.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := sourceAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - destAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= tries)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedSourceAccount, err := testRepository.GetAccount(context.Background(), sourceAccount.ID)
	require.NoError(t, err)

	updatedDestAccount, err := testRepository.GetAccount(context.Background(), destAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedSourceAccount.Balance, updatedDestAccount.Balance)

	require.Equal(t, sourceAccount.Balance-int64(tries)*amount, updatedSourceAccount.Balance)
	require.Equal(t, destAccount.Balance+int64(tries)*amount, updatedDestAccount.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	sourceAccount := createRandomAccount(t)
	destAccount := createRandomAccount(t)
	fmt.Println(">> before:", sourceAccount.Balance, destAccount.Balance)

	tries, amount := 10, int64(10)
	errs := make(chan error)

	// run tries concurent transfer transaction
	for i := 0; i < tries; i++ {
		fromAccountID := sourceAccount.ID
		toAccountID := destAccount.ID

		if i%2 == 1 {
			fromAccountID = destAccount.ID
			toAccountID = sourceAccount.ID
		}

		go func() {
			_, err := testRepository.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// checks results
	for i := 0; i < tries; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedSourceAccount, err := testRepository.GetAccount(context.Background(), sourceAccount.ID)
	require.NoError(t, err)

	updatedDestAccount, err := testRepository.GetAccount(context.Background(), destAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedSourceAccount.Balance, updatedDestAccount.Balance)

	require.Equal(t, sourceAccount.Balance, updatedSourceAccount.Balance)
	require.Equal(t, destAccount.Balance, updatedDestAccount.Balance)
}
