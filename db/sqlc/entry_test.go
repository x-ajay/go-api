package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: 1,
		Amount:    100,
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createTestEntry(t)
	})
}

func TestGetEntry(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := createTestEntry(t)
		actual, err := testQueries.GetEntry(context.Background(), expected.ID)
		require.NoError(t, err)
		require.NotEmpty(t, actual)

		require.Equal(t, expected.ID, actual.ID)
		require.Equal(t, expected.AccountID, actual.AccountID)
		require.Equal(t, expected.Amount, actual.Amount)
		require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	})

	t.Run("failure, invalid entry id", func(t *testing.T) {
		entry, err := testQueries.GetEntry(context.Background(), 0)
		require.Error(t, err)
		require.Empty(t, entry)
	})
}
