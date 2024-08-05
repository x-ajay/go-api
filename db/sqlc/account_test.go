package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	alpha    = "abcdefghijklmnopqrstuvwxyz"
	currency = []string{"USD", "EUR", "JPY"}
)

func createTestAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    generateTestOwner(),
		Balance:  generateTestBalance(),
		Currency: generateTestCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
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
	t.Run("success", func(t *testing.T) {
		createTestAccount(t)
	})

	t.Run("failure, empty owner name not allowed", func(t *testing.T) {
		arg := CreateAccountParams{
			Owner:    "",
			Balance:  generateTestBalance(),
			Currency: generateTestCurrency(),
		}
		account, err := testQueries.CreateAccount(context.Background(), arg)
		require.Error(t, err)
		require.Empty(t, account)
	})

	t.Run("failure, invalid currency should not allowed", func(t *testing.T) {
		arg := CreateAccountParams{
			Owner:    generateTestOwner(),
			Balance:  generateTestBalance(),
			Currency: "INR",
		}
		account, err := testQueries.CreateAccount(context.Background(), arg)
		require.Error(t, err)
		require.Empty(t, account)
	})

}

func TestGetAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := createTestAccount(t)
		actual, err := testQueries.GetAccount(context.Background(), expected.ID)
		require.NoError(t, err)
		require.NotEmpty(t, actual)

		require.Equal(t, expected.ID, actual.ID)
		require.Equal(t, expected.Owner, actual.Owner)
		require.Equal(t, expected.Balance, actual.Balance)
		require.Equal(t, expected.Currency, actual.Currency)
		require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	})

	t.Run("failure, invalid id", func(t *testing.T) {
		_, err := testQueries.GetAccount(context.Background(), 101010)
		require.Error(t, err)
	})
}

func TestListAccounts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			createTestAccount(t)
		}

		arg := ListAccountsParams{
			Limit:  5,
			Offset: 5,
		}
		accounts, err := testQueries.ListAccounts(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, accounts, 5)

		for _, account := range accounts {
			require.NotEmpty(t, account)
		}
	})

	t.Run("failure, invalid limit", func(t *testing.T) {
		arg := ListAccountsParams{
			Limit:  -1,
			Offset: 0,
		}
		accounts, err := testQueries.ListAccounts(context.Background(), arg)
		require.Error(t, err)
		require.Empty(t, accounts)
	})

}

func TestUpdateAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		account := createTestAccount(t)
		arg := UpdateAccountParams{
			ID:      account.ID,
			Balance: generateTestBalance(),
		}
		updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, updatedAccount)

		require.Equal(t, account.ID, updatedAccount.ID)
		require.Equal(t, account.Owner, updatedAccount.Owner)
		require.Equal(t, arg.Balance, updatedAccount.Balance)
		require.Equal(t, account.Currency, updatedAccount.Currency)
		require.Equal(t, account.CreatedAt, updatedAccount.CreatedAt)
	})

	t.Run("failure, invalid id", func(t *testing.T) {
		arg := UpdateAccountParams{
			ID:      101010,
			Balance: generateTestBalance(),
		}
		updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
		require.Error(t, err)
		require.Empty(t, updatedAccount)
	})

}

func TestDeleteAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		account := createTestAccount(t)
		err := testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)

		_, err = testQueries.GetAccount(context.Background(), account.ID)
		require.Error(t, err)
	})

	t.Run("failure, invalid id", func(t *testing.T) {
		account := createTestAccount(t)
		err := testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)

		_, err = testQueries.GetAccount(context.Background(), account.ID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
	})
}
