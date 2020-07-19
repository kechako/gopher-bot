package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Schedule struct {
	ID      int64
	Name    string
	Channel string
	Fields  string
	Command string
}

func (s *Schedule) scan(scnr scanner) error {
	err := scnr.Scan(&s.ID, &s.Name, &s.Channel, &s.Fields, &s.Command)
	if err != nil {
		return fmt.Errorf("failed to scan schedule: %w", err)
	}

	return nil
}

func (db *DB) FindSchedule(ctx context.Context, id int64) (*Schedule, error) {
	row := db.db.QueryRowContext(ctx, "select id, name, channel, fields, command from schedules where id = ?;", id)

	var s Schedule
	err := s.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find the schedule: %w", err)
	}

	return &s, nil
}

func (db *DB) FindScheduleByName(ctx context.Context, name string) (*Schedule, error) {
	row := db.db.QueryRowContext(ctx, "select id, name, channel, fields, command from schedules where name = ?;", name)

	var s Schedule
	err := s.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find the schedule: %w", err)
	}

	return &s, nil
}

func (db *DB) SearchSchedules(ctx context.Context) ([]*Schedule, error) {
	rows, err := db.db.QueryContext(ctx, "select id, name, channel, fields, command from schedules;")
	if err != nil {
		return nil, fmt.Errorf("failed to search the schedules: %w", err)
	}

	var sches []*Schedule
	for rows.Next() {
		var s Schedule
		if err := s.scan(rows); err != nil {
			return nil, fmt.Errorf("failed to search the schedules: %w", err)
		}

		sches = append(sches, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to search the schedules: %w", err)
	}

	return sches, nil
}

func (db *DB) SaveSchedule(ctx context.Context, s *Schedule) (err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = collectTransaction(tx, err)
	}()

	_, err = db.FindScheduleByName(ctx, s.Name)
	if err == nil {
		err = ErrDuplicated
		return
	} else if err != ErrNotFound {
		err = fmt.Errorf("failed to save the schedule: %w", err)
		return
	}

	if s.ID == 0 {
		err = db.insertSchedule(ctx, tx, s)
	} else {
		err = db.updateSchedule(ctx, tx, s)
	}

	return
}

func (db *DB) insertSchedule(ctx context.Context, tx *sql.Tx, s *Schedule) error {
	const stmt = `
	insert into schedules (name, channel, fields, command) values (?, ?, ?, ?);
	`
	res, err := tx.ExecContext(ctx, stmt, s.Name, s.Channel, s.Fields, s.Command)
	if err != nil {
		return fmt.Errorf("failed to insert the schedule: %w", err)
	}

	s.ID, _ = res.LastInsertId()

	return nil
}

func (db *DB) updateSchedule(ctx context.Context, tx *sql.Tx, s *Schedule) error {
	const stmt = `
	update schedules set name = ?, channel = ?, fields = ?, command = ? where id = ?;
	`
	res, err := tx.ExecContext(ctx, stmt, s.Name, s.Channel, s.Fields, s.Command, s.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update the schedule: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *DB) DeleteSchedule(ctx context.Context, id int64) error {
	const stmt = `
	delete from schedules where id = ?;
	`
	res, err := db.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete the schedule: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *DB) DeleteScheduleByName(ctx context.Context, name string) error {
	const stmt = `
	delete from schedules where name = ?;
	`
	res, err := db.db.ExecContext(ctx, stmt, name)
	if err != nil {
		return fmt.Errorf("failed to delete the schedule: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}
