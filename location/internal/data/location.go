package data

import (
	"context"

	"github.com/kechako/gopher-bot/internal/store"
	"golang.org/x/xerrors"
)

const keyPrefix = "location_"

var (
	ErrKeyNotFound = xerrors.New("key is not found")
	ErrDuplicated  = xerrors.New("location is duplicated")
)

type Location struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func GetLocation(ctx context.Context, name string) (*Location, error) {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, xerrors.New("could not get database store from context")
	}

	var loc Location
	err := s.View(func(tx *store.Tx) error {
		return tx.Get(keyPrefix+name, &loc)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, xerrors.Errorf("failed to get location of %s: %w", name, err)
	}

	return &loc, nil
}

func GetLocations(ctx context.Context) ([]*Location, error) {
	var locs []*Location

	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return nil, xerrors.New("could not get database store from context")
	}

	err := s.View(func(tx *store.Tx) error {
		it := tx.NewIterator()
		defer it.Close()

		for it.Seek(keyPrefix); it.ValidPrefix(keyPrefix); it.Next() {
			var loc Location
			_, err := it.Get(&loc)
			if err != nil {
				return xerrors.Errorf("failed to get locations: %w", err)
			}
			locs = append(locs, &loc)
		}

		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to get locations", err)
	}

	return locs, nil
}

func AddLocation(ctx context.Context, loc *Location) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return xerrors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Location
		if err := tx.Get(keyPrefix+loc.Name, &old); err != nil {
			return ErrDuplicated
		}

		return tx.Set(keyPrefix+loc.Name, loc)
	})
	if err != nil {
		if err == ErrDuplicated {
			return err
		}

		return xerrors.Errorf("failed to add location of %s: %w", loc.Name, err)
	}

	return nil
}

func UpdateLocation(ctx context.Context, loc *Location) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return xerrors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		var old Location
		if err := tx.Get(keyPrefix+loc.Name, &old); err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return tx.Set(keyPrefix+loc.Name, loc)
	})
	if err != nil {
		if err == ErrKeyNotFound {
			return err
		}

		return xerrors.Errorf("failed to update location of %s: %w", loc.Name, err)
	}

	return nil
}

func RemoveLocation(ctx context.Context, name string) error {
	s, ok := store.StoreFromContext(ctx)
	if !ok {
		return xerrors.New("could not get database store from context")
	}

	err := s.Update(func(tx *store.Tx) error {
		return tx.Delete(keyPrefix + name)
	})
	if err != nil {
		if err == store.ErrKeyNotFound {
			return ErrKeyNotFound
		}

		return xerrors.Errorf("failed to remove location of %s: %w", name, err)
	}

	return nil
}
