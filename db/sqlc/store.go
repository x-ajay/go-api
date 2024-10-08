package db

import (
	"context"
	"database/sql"
	"fmt"
)

//go:generate mockgen -source=store.go -destination=mocks/store.go -package=mock Store
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int32 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				arg.FromAccountID,
				arg.ToAccountID,
				-arg.Amount,
				arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				q,
				arg.ToAccountID,
				arg.FromAccountID,
				arg.Amount,
				-arg.Amount,
			)
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	queries *Queries,
	fromAccountID, toAccountID int64,
	fromAmount, toAmount int32,
) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{ID: fromAccountID, Amount: fromAmount})
	if err != nil {
		return fromAccount, toAccount, err
	}

	toAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{ID: toAccountID, Amount: toAmount})
	if err != nil {
		return fromAccount, toAccount, err
	}
	return fromAccount, toAccount, nil
}
