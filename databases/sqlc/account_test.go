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

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	if err != nil {
		log.Fatal("cannot create account: ", err)
	}

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
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	getAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getAccount1)

	require.Equal(t, account1.ID, getAccount1.ID)
	require.Equal(t, account1.Owner, getAccount1.Owner)
	require.Equal(t, account1.Balance, getAccount1.Balance)
	require.Equal(t, account1.Currency, getAccount1.Currency)
	require.WithinDuration(t, account1.CreatedAt, getAccount1.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: 1500,
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Owner, accountUpdated.Owner)

	require.Equal(t, arg.Balance, accountUpdated.Balance)
	require.Equal(t, account.Currency, accountUpdated.Currency)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	getAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, getAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Offset: 5,
		Limit:  5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
