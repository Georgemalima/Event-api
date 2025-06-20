package main

import (
	"net/http"

	"github.com/sikozonpc/social/internal/store"
)

type CreateCardPayload struct {
	EventID        int64  `json:"event_id" validate:"required"`
	ImagePath      string `json:"name"`
	GuestID        int64  `json:"guest_id"`
	CardTemplateID int64
}

func (app *application) createCardHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCardPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	card := &store.Card{
		EventID:        payload.EventID,
		GuestID:        payload.GuestID,
		CardTemplateID: payload.CardTemplateID,
		ImagePath:      payload.ImagePath,
	}

	ctx := r.Context()

	if err := app.store.Cards.Create(ctx, card); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, card); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
