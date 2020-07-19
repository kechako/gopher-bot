package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrDuplicated = errors.New("dupplicated")
)

type DB struct {
	db *sql.DB
}

func Open(name string) (*DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open database [%s]: %w", name, err)
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func collectTransaction(tx *sql.Tx, err error) error {
	if err == nil {
		err = tx.Commit()
	}

	if err == nil {
		return nil
	}

	if rerr := tx.Rollback(); rerr != nil {
		return fmt.Errorf("failed to rollback (with %v): %w", err, rerr)
	}

	return err
}

type scanner interface {
	Scan(dest ...interface{}) error
}

type contextKey string

var databaseContextKey contextKey = "database"

// ContextWithDB returns a context.Context including the db from the parent.
func ContextWithDB(parent context.Context, db *DB) context.Context {
	return context.WithValue(parent, databaseContextKey, db)
}

// FromContext returns a *DB from the ctx.
func FromContext(ctx context.Context) (*DB, bool) {
	db, ok := ctx.Value(databaseContextKey).(*DB)
	if !ok {
		return nil, false
	}

	return db, true
}
