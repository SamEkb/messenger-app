package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type QueryEngine interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type DB struct {
	*sqlx.DB
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

func (db *DB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return db.DB.GetContext(ctx, dest, query, args...)
}

type Tx struct {
	*sqlx.Tx
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.Tx.ExecContext(ctx, query, args...)
}

func (tx *Tx) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return tx.Tx.GetContext(ctx, dest, query, args...)
}

func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
