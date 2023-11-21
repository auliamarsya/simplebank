package databases

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/auliamarsya/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	amount := utils.RandomMoney()

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

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
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	getTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)

	require.Equal(t, transfer.ID, getTransfer.ID)
	require.Equal(t, transfer.FromAccountID, getTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, getTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, getTransfer.Amount)

	require.WithinDuration(t, transfer.CreatedAt, getTransfer.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	toAmount := utils.RandomMoney()

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: toAmount,
	}

	transferUpdated, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transferUpdated)

	require.Equal(t, transfer.ID, transferUpdated.ID)
	require.Equal(t, transfer.FromAccountID, transferUpdated.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferUpdated.ToAccountID)
	require.Equal(t, toAmount, transferUpdated.Amount)

	require.WithinDuration(t, transfer.CreatedAt, transferUpdated.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	getTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, getTransfer)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 4; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Offset: 2,
		Limit:  2,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
