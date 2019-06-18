package data

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kechako/gopher-bot/internal/store"
)

var testLocs = []*Location{
	{
		Name:      "AAAA",
		Latitude:  35.1234,
		Longitude: 138.1234,
	},
	{
		Name:      "BBBB",
		Latitude:  36.2345,
		Longitude: 139.2345,
	},
	{
		Name:      "CCCC",
		Latitude:  37.3456,
		Longitude: 140.3456,
	},
	{
		Name:      "DDDD",
		Latitude:  38.4567,
		Longitude: 141.4567,
	},
}

func makeTestDir() (string, func(), error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot-store-test")
	if err != nil {
		return "", nil, err
	}

	return dir, func() {
		os.RemoveAll(dir)
	}, nil
}

func Test_Location(t *testing.T) {
	dir, remove, err := makeTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer remove()

	s, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	ctx := store.ContextWithStore(context.Background(), s)

	_, err = GetLocation(ctx, "NOT_EXISTS")
	if err != ErrKeyNotFound {
		t.Error("GetLocation must be return ErrKeyNotFound")
	}

	for _, tt := range testLocs {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := AddLocation(ctx, tt)
			if err != nil {
				t.Error(err)
			}
		})
	}

	locs, err := GetLocations(ctx)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(locs, testLocs); diff != "" {
		t.Errorf("failed to get locations from database: (-got +want)\n%s", diff)
	}

	updateLoc := &Location{
		Name:      "NOT_EXISTS",
		Latitude:  50.9875,
		Longitude: 150.9876,
	}
	err = UpdateLocation(ctx, updateLoc)
	if err != ErrKeyNotFound {
		t.Error("GetLocation must be return ErrKeyNotFound")
	}
	updateLoc.Name = "AAAA"
	err = UpdateLocation(ctx, updateLoc)
	if err != nil {
		t.Error(err)
	}

	loc, err := GetLocation(ctx, "AAAA")
	if diff := cmp.Diff(loc, updateLoc); diff != "" {
		t.Errorf("failed to get location from database: (-got +want)\n%s", diff)
	}

	err = RemoveLocation(ctx, "NOT_EXISTS")
	if err != ErrKeyNotFound {
		t.Error("RemoveLocation must be return ErrKeyNotFound")
	}

	err = RemoveLocation(ctx, "AAAA")
	if err != nil {
		t.Error(err)
	}

	_, err = GetLocation(ctx, "AAAA")
	if err != ErrKeyNotFound {
		t.Error("GetLocation must be return ErrKeyNotFound")
	}
}
