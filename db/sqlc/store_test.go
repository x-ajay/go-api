package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	n := 5
	amount := int32(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		// verify transfer
		require.NotEmpty(t, res.Transfer)
		require.NotZero(t, res.FromEntry.ID)
		require.NotZero(t, res.ToEntry.ID)

		// verify transfer in db
		_, err = store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		// verify from account entry
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		// check if from account balance is updated
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// verify to account entry
		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// check if to account balance is updated
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2 * amount, 3 * amount, ...
	}

	// check the final balance
	fromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	toAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, account1.Balance-int32(n)*amount)
	require.Equal(t, toAccount.Balance, account2.Balance+int32(n)*amount)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	n := 10
	amount := int32(100)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			var fromAccountID = account1.ID
			var toAccountID = account2.ID
			if i%2 == 0 {
				fromAccountID = account2.ID
				toAccountID = account1.ID
			}

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

	}

	// check the final balance
	fromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	toAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, account1.Balance)
	require.Equal(t, toAccount.Balance, account2.Balance)
}
