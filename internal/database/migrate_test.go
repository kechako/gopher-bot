package database

import (
	"database/sql"
	"testing"
)

func Test_getVersion(t *testing.T) {
	path, cleanup, err := makeTestDir("test.db")
	if err != nil {
		t.Fatal("failed to create test directory: ", err)
	}

	t.Cleanup(cleanup)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction: ", err)
	}

	version, err := getVersion(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if version != 0 {
		t.Errorf("got %d, want %d", version, 0)
	}

	err = setVersion(tx, 10)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	tx.Commit()

	tx, err = db.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction: ", err)
	}

	version, err = getVersion(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if version != 10 {
		t.Errorf("got %d, want %d", version, 10)
	}

	tx.Commit()
}
