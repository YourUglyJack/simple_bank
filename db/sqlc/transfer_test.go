package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account, defalutAmount bool) Transfer {

	if account1.Balance <= 0 {
		t.Fatalf("account1 balance must be > 0, got %d", account1.Balance)
	}
	var amount int64
	if defalutAmount {
		amount = 1
	} else {
		amount = rand.Int63n(account1.Balance) + 1
	}

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

// 这里只验证了 transfer表对不对，余额不需要我去验证
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2, false)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, account1, account2, false)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.NotZero(t, transfer2.CreatedAt)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 3
	for range n {
		createRandomTransfer(t, account1, account2, true)
	}

	lstArg := ListTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         int32(n),
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), lstArg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, n)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
