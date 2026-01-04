package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/schema"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/tz"
	"github.com/moroz/homeosapiens-go/services"
	eventregistrations "github.com/moroz/homeosapiens-go/tmpl/event_registrations"
	"github.com/moroz/homeosapiens-go/types"
)

var decoder = schema.NewDecoder()

type eventRegistrationController struct {
	eventService             *services.EventService
	eventRegistrationService *services.EventRegistrationService
}

func EventRegistrationController(db queries.DBTX) *eventRegistrationController {
	return &eventRegistrationController{
		eventService:             services.NewEventService(db),
		eventRegistrationService: services.NewEventRegistrationService(db),
	}
}

func (c *eventRegistrationController) New(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	event, err := c.eventService.GetEventDetailsBySlug(r.Context(), slug)
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
	location := r.Context().Value(config.LocationContextName).(*time.Location)
	countryGuess := tz.GuessRegionByTimezone(location.String())
	params := buildRegistrationParams(user, countryGuess)
	validationErrors := make(validation.Errors)

	if err := eventregistrations.New(r.Context(), event, params, validationErrors).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}

func buildRegistrationParams(user *queries.User, location *tz.TimezoneGuess) *types.CreateEventRegistrationParams {
	var params types.CreateEventRegistrationParams

	if user != nil {
		params.GivenName = user.GivenName.String()
		params.FamilyName = user.FamilyName.String()
		params.Email = user.Email.String()

		if user.Country != nil {
			params.Country = *user.Country
		}

		if user.Profession != nil {
			params.Profession = *user.Profession
		}
	}

	if params.Country == "" && location.Found {
		params.Country = location.IsoCode
	}

	return &params
}

func (c *eventRegistrationController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, err, 400)
		return
	}

	var params types.CreateEventRegistrationParams
	if err := decoder.Decode(&params, r.PostForm); err != nil {
		handleError(w, err, 400)
		return
	}

	event, err := c.eventService.GetEventById(r.Context(), params.EventID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		handleError(w, err, 404)
		return
	}

	user := r.Context().Value(config.CurrentUserContextName).(*queries.User)
	_, err = c.eventRegistrationService.CreateEventRegistration(r.Context(), user, event, &params)
	if err == nil {
		http.Redirect(w, r, fmt.Sprintf("/events/%s", event.Slug), http.StatusFound)
		return
	}

	var validationErrors validation.Errors
	ok := errors.As(err, &validationErrors)
	if ok {
		dto, err := c.eventService.GetEventDetailsForEvent(r.Context(), event)
		if err != nil {
			handleError(w, err, 500)
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := eventregistrations.New(r.Context(), dto, &params, validationErrors).Render(w); err != nil {
			handleRenderingError(w, err)
		}
	}

	if err != nil {
		handleError(w, err, 500)
		return
	}
}
