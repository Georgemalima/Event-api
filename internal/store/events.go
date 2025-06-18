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
	CardTemplate   CardTemplate
}

type EventStore struct {
	db *sql.DB
}

func (s *EventStore) GetAllEvents(ctx context.Context, fq PaginatedFeedQuery) ([]Event, error) {
	query := `
		SELECT
			e.id, e.name, e.date, e.location, e.scanned_count,
			ct.image_path,
			u.username
		FROM events e
		LEFT JOIN card_templates ct ON ct.id = e.card_template_id
		LEFT JOIN users u ON u.id = e.user_id
		WHERE
			(e.name ILIKE '%' || $3 || '%' OR e.location ILIKE '%' || $3 || '%' OR u.username ILIKE '%' || $3 || '%')
		GROUP BY e.id, u.username
		ORDER BY e.created_at ` + fq.Sort + `
		LIMIT $1 OFFSET $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Location,
			&e.ScannedCount,
			&e.CreatedAt,
			&e.CardTemplate.ImagePath,
			&e.User.Username,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
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

func (s *EventStore) Update(ctx context.Context, event *Event) error {
	query := `
		UPDATE events
		SET name = $1, date = $2, location = $3
		WHERE id = $4
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		event.Name,
		event.Date,
		event.Location,
		event.CardTemplateID,
	).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
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
