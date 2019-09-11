package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/internal/store"
)

const keyPrefix = "location_"

var (
	ErrKeyNotFound = errors.New("key is not found")
	ErrDuplicated  = errors.New("location is duplicated")
)

type Location struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func GetLocation(ctx context.Context, name string) (*Location, error) {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, errors.New("could not get database store from context")
	}

	var loc Location
	err := s.View(func(tx *store.Tx) error {
		return tx.Get(keyPrefix+name, &loc)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to get location of %s: %w", name, err)
	}

	return &loc, nil
}

func GetLocations(ctx context.Context) ([]*Location, error) {
	var locs []*Location

	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, errors.New("could not get database store from context")
	}

	err := s.View(func(tx *store.Tx) error {
		it := tx.NewIterator()
		defer it.Close()

		for it.Seek(keyPrefix); it.ValidPrefix(keyPrefix); it.Next() {
			var loc Location
			_, err := it.Get(&loc)
			if err != nil {
				return fmt.Errorf("failed to get locations: %w", err)
			}
			locs = append(locs, &loc)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	return locs, nil
}

func AddLocation(ctx context.Context, loc *Location) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Location
		if err := tx.Get(keyPrefix+loc.Name, &old); err != nil {
			if err != store.ErrKeyNotFound {
				return err
			}
		} else {
			return ErrDuplicated
		}

		return tx.Set(keyPrefix+loc.Name, loc)
	})
	if err != nil {
		if err == ErrDuplicated {
			return err
		}

		return fmt.Errorf("failed to add location of %s: %w", loc.Name, err)
	}

	return nil
}

func UpdateLocation(ctx context.Context, loc *Location) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Location
		if err := tx.Get(keyPrefix+loc.Name, &old); err != nil {
			return err
		}

		return tx.Set(keyPrefix+loc.Name, loc)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return fmt.Errorf("failed to update location of %s: %w", loc.Name, err)
	}

	return nil
}

func RemoveLocation(ctx context.Context, name string) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return errors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Location
		if err := tx.Get(keyPrefix+name, &old); err != nil {
			return err
		}

		return tx.Delete(keyPrefix + name)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return fmt.Errorf("failed to remove location of %s: %w", name, err)
	}

	return nil
}
