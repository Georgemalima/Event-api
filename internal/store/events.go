package store

import (
	"context"
	"database/sql"
	"errors"
)

type Event struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Date           string `json:"date"`
	Location       string `json:"location"`
	ScannedCount   int64  `json:"scanned_count"`
	CardTemplateID string `json:"card_template_id"`
	UserID         int64  `json:"user_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	User           User   `json:"user"`
}

type EventStore struct {
	db *sql.DB
}

func (s *EventStore) Create(ctx context.Context, tx *sql.Tx, event *Event) error {
	query := `
		INSERT INTO events (name, date, location, user_id)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		event.Name,
		event.Date,
		event.Location,
		event.UserID,
	).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventStore) GetByID(ctx context.Context, id int64) (*Event, error) {
	query := `
		SELECT id, name, date, location, scanned_count, card_template_id, user_id, created_at,  updated_at
		FROM events
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var event Event
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Date,
		&event.Location,
		&event.ScannedCount,
		&event.CardTemplateID,
		&event.UserID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &event, nil
}

func (s *EventStore) Delete(ctx context.Context, eventID int64) error {
	query := `DELETE FROM events WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, event *Event) error {
	query := `
		UPDATE events
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		event.Name,
		event.Date,
		event.ID,
		event.CardTemplateID,
	).Scan(&event)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
