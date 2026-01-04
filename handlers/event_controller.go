package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
)

type eventController struct {
	db           queries.DBTX
	eventService *services.EventService
}

func EventController(db queries.DBTX) *eventController {
	return &eventController{db, services.NewEventService(db)}
}

func (c *eventController) Show(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	event, err := c.eventService.GetEventDetailsBySlug(r.Context(), slug)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Not found", 404)
		return
	}

	if err != nil {
		log.Printf("Error fetching event by slug %s: %s", slug, err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := events.Show(r.Context(), event).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
