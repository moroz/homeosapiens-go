package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	eventregistrations "github.com/moroz/homeosapiens-go/tmpl/event_registrations"
	"github.com/moroz/homeosapiens-go/types"
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
	user := r.Context().Value(config.CurrentUserContextName).(*queries.User)
	params := buildRegistrationParamsFromUser(user)

	if err := eventregistrations.New(r.Context(), event, params).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}

func buildRegistrationParamsFromUser(user *queries.User) *types.CreateEventRegistrationParams {
	var params types.CreateEventRegistrationParams

	if user == nil {
		return &params
	}

	params.GivenName = user.GivenName.String()
	params.FamilyName = user.FamilyName.String()
	params.Email = user.Email.String()

	if user.Country != nil {
		params.Country = *user.Country
	}
	if user.Profession != nil {
		params.Profession = *user.Profession
	}
	return &params
}
