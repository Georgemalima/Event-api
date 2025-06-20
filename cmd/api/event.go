package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sikozonpc/social/internal/store"
)

type eventKey string

const eventCtx eventKey = "event"

type CreateEventPayload struct {
	Name           string `json:"title" validate:"required,max=100"`
	Date           string `json:"date" validate:"required"`
	Location       string `json:"location"`
	CardTemplateID string `json:"card_template_id"`
}

func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateEventPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	user := getUserFromContext(r)

	event := &store.Event{
		Name:     payload.Name,
		Date:     payload.Date,
		Location: payload.Location,
		UserID:   user.ID,
	}

	ctx := r.Context()

	if err := app.store.Events.Create(ctx, event); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, event); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	events, err := app.store.Events.GetAllEvents(ctx, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, events); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getEventHandler(w http.ResponseWriter, r *http.Request) {
	event := getEventFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, event); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "eventID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Events.Delete(ctx, id); err != nil {
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

type UpdateEventPayload struct {
	Name           string `json:"name" validate:"omitempty"`
	Date           string `json:"date"`
	Location       string `json:"location"`
	CardTemplateID string `json:"card_template_id"`
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	event := getEventFromCtx(r)

	var payload UpdateEventPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name != "" {
		event.Name = payload.Name
	}
	if payload.Date != "" {
		event.Date = payload.Date
	}
	if payload.Location != "" {
		event.Location = payload.Location
	}
	if payload.CardTemplateID != "" {
		event.CardTemplateID = payload.CardTemplateID
	}

	ctx := r.Context()

	if err := app.updateEvent(ctx, event); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, event); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getEventFromCtx(r *http.Request) *store.Event {
	event, _ := r.Context().Value(eventCtx).(*store.Event)
	return event
}

func (app *application) updateEvent(ctx context.Context, event *store.Event) error {
	if err := app.store.Events.Update(ctx, event); err != nil {
		return err
	}

	app.cacheStorage.Users.Delete(ctx, event.UserID)
	return nil
}
