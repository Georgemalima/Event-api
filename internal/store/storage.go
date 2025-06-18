package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Events interface {
		GetByID(context.Context, int64) (*Event, error)
		Create(context.Context, *sql.Tx, *Event) error
		Delete(context.Context, int64) error
		Update(context.Context, *Event) error
		GetAllEvents(context.Context, PaginatedFeedQuery) ([]Event, error)
	}
	Guests interface {
		Create(ctx context.Context, tx *sql.Tx, guest *Guest) error
		Delete(ctx context.Context, guestID int64) error
		GetByID(ctx context.Context, id int64) (*Guest, error)
		GetGuests(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]Guest, error)
		Update(ctx context.Context, tx *sql.Tx, guest *Guest) error
	}
	Users interface {
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Cards interface {
		Create(ctx context.Context, tx *sql.Tx, card *Card) error
		Delete(ctx context.Context, cardID int64) error
		GetByID(ctx context.Context, id int64) (*Card, error)
		GetCards(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]Card, error)
		Update(ctx context.Context, tx *sql.Tx, card *Card) error
	}
	CardTemplates interface {
		Create(ctx context.Context, tx *sql.Tx, card *CardTemplate) error
		Delete(ctx context.Context, cardID int64) error
		GetByID(ctx context.Context, id int64) (*CardTemplate, error)
		GetCards(ctx context.Context, eventId int64, fq PaginatedFeedQuery) ([]CardTemplate, error)
		Update(ctx context.Context, tx *sql.Tx, card *CardTemplate) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Events:        &EventStore{db},
		Guests:        &GuestStore{db},
		Users:         &UserStore{db},
		Cards:         &CardStore{db},
		CardTemplates: &CardTemplateStore{db},
		Roles:         &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
