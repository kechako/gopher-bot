package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Location struct {
	ID        int64
	Name      string
	Latitude  float32
	Longitude float32
}

func (l *Location) scan(scnr scanner) error {
	err := scnr.Scan(&l.ID, &l.Name, &l.Latitude, &l.Longitude)
	if err != nil {
		return fmt.Errorf("failed to scan location: %w", err)
	}

	return nil
}

func (db *DB) FindLocation(ctx context.Context, id int64) (*Location, error) {
	row := db.db.QueryRowContext(ctx, "select id, name, latitude, longitude from locations where id = ?;", id)

	var l Location
	err := l.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find the location: %w", err)
	}

	return &l, nil
}

func (db *DB) FindLocationByName(ctx context.Context, name string) (*Location, error) {
	row := db.db.QueryRowContext(ctx, "select id, name, latitude, longitude from locations where name = ?;", name)

	var l Location
	err := l.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find the location: %w", err)
	}

	return &l, nil
}

func (db *DB) SearchLocations(ctx context.Context) ([]*Location, error) {
	rows, err := db.db.QueryContext(ctx, "select id, name, latitude, longitude from locations;")
	if err != nil {
		return nil, fmt.Errorf("failed to search the locations: %w", err)
	}

	var locs []*Location
	for rows.Next() {
		var l Location
		if err := l.scan(rows); err != nil {
			return nil, fmt.Errorf("failed to search the locations: %w", err)
		}

		locs = append(locs, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to search the locations: %w", err)
	}

	return locs, nil
}

func (db *DB) SaveLocation(ctx context.Context, l *Location) (err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = collectTransaction(tx, err)
	}()

	_, err = db.FindLocationByName(ctx, l.Name)
	if err == nil {
		err = ErrDuplicated
		return
	} else if err != ErrNotFound {
		err = fmt.Errorf("failed to save the location: %w", err)
		return
	}

	if l.ID == 0 {
		err = db.insertLocation(ctx, tx, l)
	} else {
		err = db.updateLocation(ctx, tx, l)
	}

	return
}

func (db *DB) insertLocation(ctx context.Context, tx *sql.Tx, l *Location) error {
	const stmt = `
	insert into locations (name, latitude, longitude) values (?, ?, ?);
	`
	res, err := tx.ExecContext(ctx, stmt, l.Name, l.Latitude, l.Longitude)
	if err != nil {
		return fmt.Errorf("failed to insert the location: %w", err)
	}

	l.ID, _ = res.LastInsertId()

	return nil
}

func (db *DB) updateLocation(ctx context.Context, tx *sql.Tx, l *Location) error {
	const stmt = `
	update locations set name = ?, latitude = ?, longitude = ? where id = ?;
	`
	res, err := tx.ExecContext(ctx, stmt, l.Name, l.Latitude, l.Longitude, l.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update the location: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *DB) DeleteLocation(ctx context.Context, id int64) error {
	const stmt = `
	delete from locations where id = ?;
	`
	res, err := db.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete the location: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *DB) DeleteLocationByName(ctx context.Context, name string) error {
	const stmt = `
	delete from locations where name = ?;
	`
	res, err := db.db.ExecContext(ctx, stmt, name)
	if err != nil {
		return fmt.Errorf("failed to delete the location: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}
