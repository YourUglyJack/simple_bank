package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store 
func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}


// execTx executes a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txQ := New(tx)
	err = fn(txQ)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer 	Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry 	Entry `json:"from_entry"`  // -amount
	ToEntry 	Entry `json:"to_entry"`  // +amount
}

// TransferTx performs a money transfer from one account to the another
// It creates a transfer record, add account entries, and update account's balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// 1. create a transfer
		var err error
		// PostgreSQL需要：检查账户1和2是否存在 → 对accounts.id=1和id=2加共享锁
		// 		-- ✅ 此时：共享锁可以共存，两个INSERT都成功
		// -- 账户1上有：S_A1（事务A的共享锁）+ S_B1（事务B的共享锁）
		// -- 账户2上有：S_A2（事务A的共享锁）+ S_B2（事务B的共享锁）
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3. update account's amount
		// NOTE: 不加锁查询的话，多个并发goroutine可能在开始的时候拿到相同的balance，然后进行加减操作，这样会出现数据混乱
		// NOTE： FOR NO KEY UPDATE
		// 事务A 尝试获取排他锁 X_A1，失败，因为要等 事务B 释放 账户1上的 Shared 锁, S_B1
		// 事务B 尝试获取排他锁 X_B1，失败，因为要等 事务A 释放 账户1上的 Shared 锁, A_B1
		// 互相等待，所以死锁
		// 因此 origin tx 是错误的，不能 FOR UPDATE 要 FOR NO KEY UPDATE
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID: arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }


		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID: arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}