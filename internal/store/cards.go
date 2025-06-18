package store

import (
	"context"
	"database/sql"
	"errors"
)

type Card struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	EventID   int64  `json:"event_id"`
	GuestID   int64  `json:"guest_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Guest     *Guest `json:"guest"`
	Event     Event  `json:"event"`
}
type CardStore struct {
	db *sql.DB
}

func (s *CardStore) GetCards(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]Card, error) {
	query := `
		SELECT
			c.id, c.image_path,
			gs.name, gs.phone_number
		FROM cards c
		LEFT JOIN guests gs ON gs.id = c.guest_id
		WHERE
			(gs.name ILIKE '%' || $4 || '%' OR gs.phone_number ILIKE '%' || $4 || '%')
		GROUP BY c.id, gs.name
		ORDER BY c.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		err := rows.Scan(
			&card.ID,
			&card.ImagePath,
			&card.Guest.Name,
			&card.Guest.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (s *CardStore) GetByID(ctx context.Context, id int64) (*Card, error) {
	query := `
		SELECT id, image_path, guest_id, created_at,  updated_at,
			gs.name, gs.phone_number
		FROM cards c
		LEFT JOIN guests gs ON gs.id = c.guest_id
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var card Card
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.ImagePath,
		&card.GuestID,
		&card.CreatedAt,
		&card.UpdatedAt,
		&card.Guest.Name,
		&card.Guest.PhoneNumber,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &card, nil
}

func (s *CardStore) Create(ctx context.Context, tx *sql.Tx, card *Card) error {
	query := `
		INSERT INTO cards (event_id, guest_id, image_path)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		card.EventID,
		card.GuestID,
		card.ImagePath,
	).Scan(
		&card.ID,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CardStore) Delete(ctx context.Context, cardID int64) error {
	query := `DELETE FROM cards WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, cardID)
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

func (s *CardStore) Update(ctx context.Context, tx *sql.Tx, card *Card) error {
	query := `
		UPDATE cards
		SET event_id = $1, guest_id = $2, image_path = $3
		WHERE id = $4
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		card.EventID,
		card.GuestID,
		card.ImagePath,
	).Scan(
		&card.ID,
		&card.CreatedAt,
		&card.UpdatedAt,
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
