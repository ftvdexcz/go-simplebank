package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ftvdexcz/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func createRandomAccount(t *testing.T) Account{
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestGetAccount(t *testing.T){
	// create account 
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)  

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	assert.Equal(t, acc1.ID, acc2.ID)
	assert.Equal(t, acc1.Owner, acc2.Owner)
	assert.Equal(t, acc1.Balance, acc2.Balance)
	assert.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestUpdateAccount(t *testing.T){
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	assert.Equal(t, acc1.ID, acc2.ID)
	assert.Equal(t, acc1.Owner, acc2.Owner)
	assert.Equal(t, arg.Balance, acc2.Balance)
	assert.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T){
	acc1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)

	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T){
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts{
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}