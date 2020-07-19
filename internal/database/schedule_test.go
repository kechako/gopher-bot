package database

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testSchedules = []*Schedule{
	{
		Name:    "AAAA",
		Channel: "#test1",
		Fields:  "0 9-15 * * 1-5",
		Command: "aaaaaa",
	},
	{
		Name:    "BBBB",
		Channel: "#test2",
		Fields:  "0 10-16 * * 1-5",
		Command: "bbbbbb",
	},
	{
		Name:    "CCCC",
		Channel: "#test3",
		Fields:  "0 11-17 * * 1-5",
		Command: "cccccc",
	},
	{
		Name:    "DDDD",
		Channel: "#test4",
		Fields:  "0 12-18 * * 1-5",
		Command: "dddddd",
	},
}

func Test_Schedule(t *testing.T) {
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

	for _, tt := range testSchedules {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := db.SaveSchedule(ctx, tt)
			if err != nil {
				t.Error(err)
			}
		})
	}

	locs, err := db.SearchSchedules(ctx)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(locs, testSchedules); diff != "" {
		t.Errorf("failed to get schedules from database: (-got +want)\n%s", diff)
	}

	_, err = db.FindSchedule(ctx, -1 /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.FindSchedule must be return ErrNotFound, got %v", err)
	}
	_, err = db.FindScheduleByName(ctx, "NOT_EXISTS" /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.FindScheduleByName must be return ErrNotFound, got %v", err)
	}

	updateSch := &Schedule{
		ID:      -1, // the key does not exist
		Name:    "EEEE",
		Channel: "#test100",
		Fields:  "0 15-21 * * 1-5",
		Command: "xxxxxx",
	}
	err = db.SaveSchedule(ctx, updateSch)
	if err != ErrNotFound {
		t.Errorf("DB.SaveSchedule must be return ErrNotFound, got %v", err)
	}
	updateSch.ID = 1
	err = db.SaveSchedule(ctx, updateSch)
	if err != nil {
		t.Error(err)
	}

	loc, err := db.FindSchedule(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(loc, updateSch); diff != "" {
		t.Errorf("failed to get schedule from database: (-got +want)\n%s", diff)
	}
	loc, err = db.FindScheduleByName(ctx, "EEEE")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(loc, updateSch); diff != "" {
		t.Errorf("failed to get schedule from database: (-got +want)\n%s", diff)
	}

	dupSch := &Schedule{
		Name:    "BBBB",
		Channel: "#test100",
		Fields:  "0 15-21 * * 1-5",
		Command: "xxxxxx",
	}
	err = db.SaveSchedule(ctx, dupSch)
	if err != ErrDuplicated {
		t.Errorf("DB.SaveSchedule must be return ErrDuplicated, got %v", err)
	}

	err = db.DeleteSchedule(ctx, -1 /* the key does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.DeleteSchedule must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteSchedule(ctx, 1)
	if err != nil {
		t.Error(err)
	}

	_, err = db.FindSchedule(ctx, 1)
	if err != ErrNotFound {
		t.Errorf("DB.FindSchedule must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteScheduleByName(ctx, "NOT_EXISTS" /* the name does not exist */)
	if err != ErrNotFound {
		t.Errorf("DB.DeleteScheduleByName must be return ErrNotFound, got %v", err)
	}

	err = db.DeleteScheduleByName(ctx, "BBBB")
	if err != nil {
		t.Error(err)
	}

	_, err = db.FindScheduleByName(ctx, "BBBB")
	if err != ErrNotFound {
		t.Errorf("DB.FindScheduleByName must be return ErrNotFound, got %v", err)
	}
}
