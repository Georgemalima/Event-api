package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sikozonpc/social/internal/store"
)

type guestKey string

const guestCtx guestKey = "guest"

type CreateEventGuests struct {
}

func (app *application) getEventGuestsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "eventID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	fq, err = fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	guests, err := app.store.Guests.GetGuests(ctx, id, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, guests); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteGuestHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "guestID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Guests.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
