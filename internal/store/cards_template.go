package store

import (
	"context"
	"database/sql"
	"errors"
)

type CardTemplate struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CardTemplateStore struct {
	db *sql.DB
}

func (s *CardTemplateStore) GetCards(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]CardTemplate, error) {
	query := `
		SELECT
			ct.id, ct.image_path
		FROM card_templates ct
		ORDER BY ct.created_at ` + fq.Sort + `
		LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cards []CardTemplate
	for rows.Next() {
		var card CardTemplate
		err := rows.Scan(
			&card.ID,
			&card.ImagePath,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (s *CardTemplateStore) GetByID(ctx context.Context, id int64) (*CardTemplate, error) {
	query := `
		SELECT id, image_path, created_at,  updated_at
		FROM card_templates
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var card CardTemplate
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.ImagePath,
		&card.CreatedAt,
		&card.UpdatedAt,
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

func (s *CardTemplateStore) Create(ctx context.Context, tx *sql.Tx, card *CardTemplate) error {
	query := `
		INSERT INTO card_templates (image_path)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
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

func (s *CardTemplateStore) Delete(ctx context.Context, cardID int64) error {
	query := `DELETE FROM card_templates WHERE id = $1`

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

func (s *CardTemplateStore) Update(ctx context.Context, tx *sql.Tx, card *CardTemplate) error {
	query := `
		UPDATE card_templates
		SET image_path = $1
		WHERE id = $2
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
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
