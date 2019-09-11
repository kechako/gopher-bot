package data

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kechako/gopher-bot/internal/store"
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

func makeTestDir() (string, func(), error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot-store-test")
	if err != nil {
		return "", nil, err
	}

	return dir, func() {
		os.RemoveAll(dir)
	}, nil
}

func Test_Schedule(t *testing.T) {
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

	_, err = GetSchedule(ctx, "NOT_EXISTS")
	if err != ErrKeyNotFound {
		t.Error("GetSchedule must be return ErrKeyNotFound")
	}

	for _, tt := range testSchedules {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			err := AddSchedule(ctx, tt)
			if err != nil {
				t.Error(err)
			}
		})
	}

	sches, err := GetSchedules(ctx)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(sches, testSchedules); diff != "" {
		t.Errorf("failed to get schedules from database: (-got +want)\n%s", diff)
	}

	updateSch := &Schedule{
		Name:    "NOT_EXISTS",
		Channel: "#test100",
		Fields:  "0 15-21 * * 1-5",
		Command: "xxxxxx",
	}
	err = UpdateSchedule(ctx, updateSch)
	if err != ErrKeyNotFound {
		t.Error("GetSchedule must be return ErrKeyNotFound")
	}
	updateSch.Name = "AAAA"
	err = UpdateSchedule(ctx, updateSch)
	if err != nil {
		t.Error(err)
	}

	sch, err := GetSchedule(ctx, "AAAA")
	if diff := cmp.Diff(sch, updateSch); diff != "" {
		t.Errorf("failed to get schedule from database: (-got +want)\n%s", diff)
	}

	err = RemoveSchedule(ctx, "NOT_EXISTS")
	if err != ErrKeyNotFound {
		t.Error("RemoveSchedule must be return ErrKeyNotFound")
	}

	err = RemoveSchedule(ctx, "AAAA")
	if err != nil {
		t.Error(err)
	}

	_, err = GetSchedule(ctx, "AAAA")
	if err != ErrKeyNotFound {
		t.Error("GetSchedule must be return ErrKeyNotFound")
	}
}
