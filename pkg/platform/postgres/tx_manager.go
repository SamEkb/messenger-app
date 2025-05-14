package postgres

import (
	"context"
)

type txKeyType struct{}

var txKey txKeyType

type TxManager struct {
	db *DB
}

func NewTxManager(db *DB) *TxManager {
	return &TxManager{db: db}
}

func (tm *TxManager) RunTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	txWrap := &Tx{Tx: tx}
	txCtx := context.WithValue(ctx, txKey, txWrap)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tm *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	if tx, ok := ctx.Value(txKey).(*Tx); ok {
		return tx
	}
	return tm.db
}
