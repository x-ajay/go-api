package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T) (Account, Account, Transfer) {

	// Create a from/to account
	createAccountParam := CreateAccountParams{
		Owner:    generateTestOwner(),
		Balance:  generateTestBalance(),
		Currency: generateTestCurrency(),
	}
	fromAccount, err := testQueries.CreateAccount(context.Background(), createAccountParam)
	require.NoError(t, err)

	createAccountParam = CreateAccountParams{
		Owner:    generateTestOwner(),
		Balance:  generateTestBalance(),
		Currency: generateTestCurrency(),
	}
	toAccount, err := testQueries.CreateAccount(context.Background(), createAccountParam)
	require.NoError(t, err)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        100,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return fromAccount, toAccount, transfer
}

func TestCreateTransfer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createTestTransfer(t)
	})
}

func TestGetTransfer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, _, transfer := createTestTransfer(t)
		transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer2)

		require.Equal(t, transfer.ID, transfer2.ID)
		require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
		require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	})

	t.Run("failure", func(t *testing.T) {
		transfer, err := testQueries.GetTransfer(context.Background(), 999)
		require.Error(t, err)
		require.Empty(t, transfer)
	})
}
