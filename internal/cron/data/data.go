// Package data accesses stored cron data.
package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/internal/store"
)

const keyPrefix = "cron_"

func genKey(name string) string {
	return keyPrefix + name
}

var (
	ErrKeyNotFound = errors.New("key is not found")
	ErrDuplicated  = errors.New("schedule is duplicated")
)

type Schedule struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Fields  string `json:"fields"`
	Command string `json:"command"`
}

func GetSchedule(ctx context.Context, name string) (*Schedule, error) {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, errors.New("could not get database store from context")
	}

	var sch Schedule
	err := s.View(func(tx *store.Tx) error {
		return tx.Get(genKey(name), &sch)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to get schedule of %s: %w", name, err)
	}

	return &sch, nil
}

func GetSchedules(ctx context.Context) ([]*Schedule, error) {
	var sches []*Schedule

	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, errors.New("could not get database store from context")
	}

	err := s.View(func(tx *store.Tx) error {
		it := tx.NewIterator()
		defer it.Close()

		for it.Seek(keyPrefix); it.ValidPrefix(keyPrefix); it.Next() {
			var sch Schedule
			_, err := it.Get(&sch)
			if err != nil {
				return fmt.Errorf("failed to get schedule: %w", err)
			}
			sches = append(sches, &sch)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return sches, nil
}

func AddSchedule(ctx context.Context, sch *Schedule) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Schedule
		if err := tx.Get(genKey(sch.Name), &old); err != nil {
			if err != store.ErrKeyNotFound {
				return err
			}
		} else {
			return ErrDuplicated
		}

		return tx.Set(genKey(sch.Name), sch)
	})
	if err != nil {
		if err == ErrDuplicated {
			return err
		}

		return fmt.Errorf("failed to add schedule of %s: %w", sch.Name, err)
	}

	return nil
}

func UpdateSchedule(ctx context.Context, sch *Schedule) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Schedule
		if err := tx.Get(genKey(sch.Name), &old); err != nil {
			return err
		}

		return tx.Set(genKey(sch.Name), sch)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return fmt.Errorf("failed to update schedule of %s: %w", sch.Name, err)
	}

	return nil
}

func RemoveSchedule(ctx context.Context, name string) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Schedule
		if err := tx.Get(genKey(name), &old); err != nil {
			return err
		}

		return tx.Delete(genKey(name))
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return fmt.Errorf("failed to remove schedule of %s: %w", name, err)
	}

	return nil
}
