package databases

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := CreateRandomAccount(t)
	accountTo := CreateRandomAccount(t)
	fmt.Println(">>before:", accountFrom.Balance, accountTo.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountFrom.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountFrom.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTo.ID, toAccount.ID)

		fmt.Println(">>tx:", fromAccount.Balance, toAccount.Balance)

		diffBalanceFrom := accountFrom.Balance - fromAccount.Balance
		diffBalanceTo := toAccount.Balance - accountTo.Balance
		require.Equal(t, diffBalanceFrom, diffBalanceTo)
		require.True(t, diffBalanceFrom > 0)
		require.True(t, diffBalanceFrom%amount == 0)

		key := int(diffBalanceFrom / amount)
		require.True(t, key >= 1 && key <= n)
		require.NotContains(t, existed, key)
		existed[key] = true
	}

	updatedFrom, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updatedTo, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	fmt.Println(">>after:", updatedFrom.Balance, updatedTo.Balance)

	require.Equal(t, accountFrom.Balance-int64(n)*amount, updatedFrom.Balance)
	require.Equal(t, accountTo.Balance+int64(n)*amount, updatedTo.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := CreateRandomAccount(t)
	accountTo := CreateRandomAccount(t)
	fmt.Println(">>before:", accountFrom.Balance, accountTo.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountFrom.ID
		toAccountID := accountTo.ID

		if i % 2 == 1 {
			fromAccountID = accountTo.ID
			toAccountID = accountFrom.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedFrom, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updatedTo, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	fmt.Println(">>after:", updatedFrom.Balance, updatedTo.Balance)

	require.Equal(t, accountFrom.Balance, updatedFrom.Balance)
	require.Equal(t, accountTo.Balance, updatedTo.Balance)
}
