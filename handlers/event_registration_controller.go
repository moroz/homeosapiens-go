package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/tz"
	"github.com/moroz/homeosapiens-go/services"
	eventregistrations "github.com/moroz/homeosapiens-go/tmpl/event_registrations"
	"github.com/moroz/homeosapiens-go/types"
)

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

func (c *eventRegistrationController) New(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)
	slug := r.Param("slug")
	event, err := c.eventService.GetEventDetailsBySlug(r.Request().Context(), slug, ctx.User)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("Event with the slug: %s not found", slug)
		return echo.ErrNotFound
	}
	if err != nil {
		log.Printf("Error when fetching event: %s", err)
		return err
	}
	location := ctx.Timezone
	countryGuess := tz.GuessRegionByTimezone(location.String())
	params := buildRegistrationParams(ctx.User, countryGuess)
	validationErrors := make(validation.Errors)

	return eventregistrations.New(ctx, event, params, validationErrors).Render(r.Response())
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

func (c *eventRegistrationController) Create(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)

	var params types.CreateEventRegistrationParams
	if err := r.Bind(&params); err != nil {
		return echo.ErrBadRequest
	}

	event, err := c.eventService.GetEventById(r.Request().Context(), params.EventID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	user := ctx.User
	_, err = c.eventRegistrationService.CreateEventRegistration(r.Request().Context(), user, event, &params)
	if err == nil {
		return r.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
	}

	var validationErrors validation.Errors
	ok := errors.As(err, &validationErrors)
	if ok {
		dto, err := c.eventService.GetEventDetailsForEvent(r.Request().Context(), event, user)
		if err != nil {
			return err
		}

		r.Response().WriteHeader(http.StatusUnprocessableEntity)
		return eventregistrations.New(ctx, dto, &params, validationErrors).Render(r.Response())
	}

	return err
}
