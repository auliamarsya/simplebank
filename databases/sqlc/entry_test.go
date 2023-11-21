package databases

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/auliamarsya/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	if err != nil {
		log.Fatal("cannot create entry: ", err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	getEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getEntry)

	require.Equal(t, entry.ID, getEntry.ID)
	require.Equal(t, entry.AccountID, getEntry.AccountID)
	require.Equal(t, entry.Amount, getEntry.Amount)

	require.WithinDuration(t, entry.CreatedAt, getEntry.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)

	toAmount := utils.RandomMoney()

	arg := UpdateEntryParams{
		ID:     entry.ID,
		Amount: toAmount,
	}

	entryUpdated, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entryUpdated)

	require.Equal(t, entry.ID, entryUpdated.ID)
	require.Equal(t, entry.AccountID, entryUpdated.AccountID)
	require.Equal(t, arg.Amount, entryUpdated.Amount)

	require.WithinDuration(t, entry.CreatedAt, entryUpdated.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	entryDeleted, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, entryDeleted)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 6; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  3,
		Offset: 3,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 3)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
