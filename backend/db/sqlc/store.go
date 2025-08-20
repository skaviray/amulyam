package db

import (
	"database/sql"
)

type Store interface {
	Querier
	// TransferTX(ctx context.Context, req TransferTXParams) (TransferTXResponse, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// type TransferTXParams struct {
// 	Amount      int64 `json:"amount"`
// 	FromAccount int64 `json:"from_account"`
// 	ToAccount   int64 `json:"to_account"`
// }

// type TransferTXResponse struct {
// 	Transfer    Transfer `json:"transfer"`
// 	FromAccount Account  `json:"from_account"`
// 	ToAccount   Account  `json:"to_account"`
// 	FromEntry   Entry    `json:"from_entry"`
// 	ToEntry     Entry    `json:"to_entry"`
// }

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// var TxKey = struct{}{}

// type TxKey string

// func (store *SQLStore) executeTx(ctx context.Context, fn func(db *Queries) error) error {
// 	tx, err := store.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	query := New(tx)
// 	err = fn(query)
// 	if err != nil {
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx err :%v, rb err: %v", err.Error(), rbErr.Error())
// 		}
// 		return err
// 	}
// 	return tx.Commit()
// }

// func (store *SQLStore) TransferTX(ctx context.Context, req TransferTXParams) (TransferTXResponse, error) {
// 	var response TransferTXResponse
// 	err := store.executeTx(context.Background(), func(q *Queries) error {
// 		var err error
// 		// txName := ctx.Value(TxKey("txname"))
// 		createTransferParams := CreateTransferParams{
// 			FromAccount: req.FromAccount,
// 			ToAccount:   req.ToAccount,
// 			Amount:      req.Amount,
// 		}
// 		// fmt.Println(txName, "create transfer")
// 		response.Transfer, err = q.CreateTransfer(ctx, createTransferParams)
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "create From Entry")
// 		response.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
// 			AccountID: req.FromAccount,
// 			Amount:    -req.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "create To Entry")
// 		response.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
// 			AccountID: req.ToAccount,
// 			Amount:    req.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "Get from Account for update")
// 		account1, err := q.GetAccountForUpdate(ctx, req.FromAccount)
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "Update from Account balance")
// 		response.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
// 			ID:      req.FromAccount,
// 			Balance: account1.Balance - req.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "Get to Account for update")
// 		account2, err := q.GetAccountForUpdate(ctx, req.ToAccount)
// 		if err != nil {
// 			return err
// 		}
// 		// fmt.Println(txName, "Update to Account balance")
// 		response.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
// 			ID:      req.ToAccount,
// 			Balance: account2.Balance + req.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	return response, err
// }
