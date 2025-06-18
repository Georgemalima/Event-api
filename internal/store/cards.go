package store

import (
	"context"
	"database/sql"
)

type Card struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	GuestID   int64  `json:"guest_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Guest     Guest  `json:"guest"`
}

type CardStore struct {
	db *sql.DB
}

func (s *GuestStore) GetCards(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]Card, error) {
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

	}
}
