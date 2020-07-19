package database

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func Test_Location(t *testing.T) {
	path, cleanup, err := makeTestDir("test.db")
	if err != nil {
		t.Fatal("failed to create test directory: ", err)
	}

	t.Cleanup(cleanup)

	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	ctx := context.Background()

	for _, tt := range testLocs {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := db.SaveLocation(ctx, tt)
			if err != nil {
				t.Error(err)
			}
		})
	}

	locs, err := db.SearchLocations(ctx)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(locs, testLocs); diff != "" {
		t.Errorf("failed to get locations from database: (-got +want)\n%s", diff)
	}

	_, err = db.FindLocation(ctx, -1 /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.FindLocation must be return ErrNotFound, got %v", err)
	}
	_, err = db.FindLocationByName(ctx, "NOT_EXISTS" /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.FindLocationByName must be return ErrNotFound, got %v", err)
	}

	updateLoc := &Location{
		ID:        -1, // the key does not exist
		Name:      "EEEE",
		Latitude:  50.9875,
		Longitude: 150.9876,
	}
	err = db.SaveLocation(ctx, updateLoc)
	if err != ErrNotFound {
		t.Errorf("DB.SaveLocation must be return ErrNotFound, got %v", err)
	}
	updateLoc.ID = 1
	err = db.SaveLocation(ctx, updateLoc)
	if err != nil {
		t.Error(err)
	}

	loc, err := db.FindLocation(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(loc, updateLoc); diff != "" {
		t.Errorf("failed to get location from database: (-got +want)\n%s", diff)
	}
	loc, err = db.FindLocationByName(ctx, "EEEE")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(loc, updateLoc); diff != "" {
		t.Errorf("failed to get location from database: (-got +want)\n%s", diff)
	}

	dupLoc := &Location{
		Name:      "BBBB",
		Latitude:  50.9875,
		Longitude: 150.9876,
	}
	err = db.SaveLocation(ctx, dupLoc)
	if err != ErrDuplicated {
		t.Errorf("DB.SaveLocation must be return ErrDuplicated, got %v", err)
	}

	err = db.DeleteLocation(ctx, -1 /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.DeleteLocation must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteLocation(ctx, 1)
	if err != nil {
		t.Error(err)
	}

	_, err = db.FindLocation(ctx, 1)
	if err != ErrNotFound {
		t.Errorf("DB.FindLocation must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteLocationByName(ctx, "NOT_EXISTS" /* the name does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.DeleteLocationByName must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteLocationByName(ctx, "BBBB")
	if err != nil {
		t.Error(err)
	}

	_, err = db.FindLocationByName(ctx, "BBBB")
	if err != ErrNotFound {
		t.Errorf("DB.FindLocationByName must be return ErrNotFound, got %v", err)
	}
}
