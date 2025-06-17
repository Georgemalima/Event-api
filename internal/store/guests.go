package store

import (
	"context"
	"database/sql"
	"errors"
)

type Guest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CardID      int64  `json:"card_id"`
	EventID     int64  `json:"event_id"`
	Event       Event  `json:"event"`
}

type GuestStore struct {
	db *sql.DB
}

func (s *GuestStore) GetGuests(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]Guest, error) {
	// to implement
	return nil, nil
}

func (s *GuestStore) GetByID(ctx context.Context, id int64, fq PaginatedFeedQuery) (*Guest, error) {
	query := `
		SELECT id, name, email, phone_number, status, type, card_id, event_id, created_at,  updated_at
		FROM guests
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var guest Guest
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&guest.ID,
		&guest.Name,
		&guest.Email,
		&guest.PhoneNumber,
		&guest.Status,
		&guest.Type,
		&guest.CardID,
		&guest.EventID,
		&guest.CreatedAt,
		&guest.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &guest, nil
}

func (s *GuestStore) Create(ctx context.Context, tx *sql.Tx, guest *Guest) error {
	query := `
		INSERT INTO guests (name, email, phone_number, status, type, event_id)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		guest.Name,
		guest.Email,
		guest.PhoneNumber,
		guest.Status,
		guest.Type,
		guest.EventID,
	).Scan(
		&guest.ID,
		&guest.CreatedAt,
		&guest.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *GuestStore) Delete(ctx context.Context, guestID int64) error {
	query := `DELETE FROM guests WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, guestID)
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

func (s *GuestStore) Update(ctx context.Context, tx *sql.Tx, guest *Guest) error {
	query := `
		UPDATE guests
		SET name = $1, email = $2, phone_number = $3, status = $4, type = $5
		WHERE id = $6
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		guest.Name,
		guest.Email,
		guest.PhoneNumber,
		guest.Status,
		guest.Type,
	).Scan(
		&guest.ID,
		&guest.CreatedAt,
		&guest.UpdatedAt,
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
