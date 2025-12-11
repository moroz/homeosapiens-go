package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	eventregistrations "github.com/moroz/homeosapiens-go/tmpl/event_registrations"
)

type eventRegistrationController struct {
	*services.EventService
}

func EventRegistrationController(db queries.DBTX) *eventRegistrationController {
	return &eventRegistrationController{
		EventService: services.NewEventService(db),
	}
}

func (c *eventRegistrationController) New(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	event, err := c.EventService.GetEventBySlug(r.Context(), slug)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("Event with the slug: %s not found", slug)
		http.Error(w, "Not Found", 404)
		return
	}
	if err != nil {
		log.Printf("Error when fetching event: %s", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := eventregistrations.New(r.Context(), event).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
