package database

import (
	"database/sql"
	"fmt"
)

const CurrentVersion = 1

func migrate(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = collectTransaction(tx, err)
	}()

	version, err := getVersion(tx)
	if err != nil {
		return
	}

	if version == 0 {
		// new database
		err = createTable(tx)
		if err != nil {
			return
		}
	} else {
		for i := version; i < CurrentVersion; i++ {
			err = updateTable(tx, i, i+1)
			if err != nil {
				return
			}
		}
	}

	version = CurrentVersion

	err = setVersion(tx, version)
	if err != nil {
		return
	}

	return
}

func createTable(tx *sql.Tx) error {
	// locations
	locStmt := `
	create table locations (
		id        integer primary key,
		name      text unique,
		latitude  real,
		longitude real
	);
	`
	_, err := tx.Exec(locStmt)
	if err != nil {
		return fmt.Errorf("failed to create table [locations]: %w", err)
	}

	// schedules
	schStmt := `
	create table schedules (
		id      integer primary key,
		name    text unique,
		channel text,
		Fields  text,
		Command text
	);
	`
	_, err = tx.Exec(schStmt)
	if err != nil {
		return fmt.Errorf("failed to create table [schedules]: %w", err)
	}

	return nil
}

func updateTable(tx *sql.Tx, oldVersion, newVersion int) error {
	return nil
}

func getVersion(tx *sql.Tx) (version int, err error) {
	rows, err := tx.Query("PRAGMA user_version;")
	if err != nil {
		err = fmt.Errorf("failed to get user_version: %w", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			err = fmt.Errorf("failed to get user_version: %w", err)
			return
		}
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("failed to get user_version: %w", err)
		return
	}

	return
}

func setVersion(tx *sql.Tx, version int) error {
	stmt := fmt.Sprintf("PRAGMA user_version = %d;", version)
	_, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to set user_version: %w", err)
	}

	return nil
}
