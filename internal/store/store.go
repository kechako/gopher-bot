package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
)

var ErrKeyNotFound = errors.New("key not found")

type Store struct {
	db *badger.DB
}

func New(dir string) (*Store, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open databsae: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Update(fn func(tx *Tx) error) error {
	var fnErr error
	err := s.db.Update(func(txn *badger.Txn) error {
		fnErr = fn(&Tx{txn: txn})
		return fnErr
	})
	if fnErr != nil {
		return fnErr
	}
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

func (s *Store) View(fn func(tx *Tx) error) error {
	var fnErr error
	err := s.db.View(func(txn *badger.Txn) error {
		fnErr = fn(&Tx{txn: txn})
		return fnErr
	})
	if fnErr != nil {
		return fnErr
	}
	if err != nil {
		return fmt.Errorf("failed to view database: %w", err)
	}

	return nil
}

type Tx struct {
	txn *badger.Txn
}

func (tx *Tx) Set(key string, value interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("failed to encode the value: %w", err)
	}

	if err := tx.txn.Set([]byte(key), buf.Bytes()); err != nil {
		return fmt.Errorf("failed to set the value: %w", err)
	}

	return nil
}

func (tx *Tx) Get(key string, value interface{}) error {
	item, err := tx.txn.Get([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrKeyNotFound
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	buf, err := item.Value()
	if err != nil {
		return fmt.Errorf("failed to get item value: %w", err)
	}

	if err := json.Unmarshal(buf, value); err != nil {
		return fmt.Errorf("failed to decode the value: %w", err)
	}

	return nil
}

func (tx *Tx) Delete(key string) error {
	err := tx.txn.Delete([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrKeyNotFound
		}
		return fmt.Errorf("failed to delete value: %w", err)
	}

	return nil
}

func (tx *Tx) NewIterator() *Iterator {
	return &Iterator{it: tx.txn.NewIterator(badger.DefaultIteratorOptions)}
}

type Iterator struct {
	it *badger.Iterator
}

func (it *Iterator) Close() {
	it.it.Close()
}

func (it *Iterator) Seek(prefix string) {
	it.it.Seek([]byte(prefix))
}

func (it *Iterator) ValidPrefix(prefix string) bool {
	return it.it.ValidForPrefix([]byte(prefix))
}

func (it *Iterator) Next() {
	it.it.Next()
}

func (it *Iterator) Get(value interface{}) (key string, err error) {
	item := it.it.Item()
	key = string(item.Key())

	if buf, verr := item.Value(); verr != nil {
		err = fmt.Errorf("failed to get item value: %w", verr)
	} else {
		if uerr := json.Unmarshal(buf, value); uerr != nil {
			err = fmt.Errorf("failed to decode the value: %w", uerr)
		}
	}

	return
}

type contextKey string

var storeContextKey contextKey = "store"

func ContextWithStore(parent context.Context, store *Store) context.Context {
	return context.WithValue(parent, storeContextKey, store)
}

func StoreFromContext(ctx context.Context) (*Store, bool) {
	store, ok := ctx.Value(storeContextKey).(*Store)
	if !ok {
		return nil, false
	}

	return store, true
}
