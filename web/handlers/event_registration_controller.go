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

func (cc *eventRegistrationController) New(c *echo.Context) error {
	ctx := c.Get("context").(*types.CustomContext)
	slug := c.Param("slug")
	event, err := cc.eventService.GetEventDetailsBySlug(c.Request().Context(), slug, ctx.User)
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

	return eventregistrations.New(ctx, event, params, validationErrors).Render(c.Response())
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

func (cc *eventRegistrationController) Create(c *echo.Context) error {
	ctx := c.Get("context").(*types.CustomContext)

	var params types.CreateEventRegistrationParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest
	}

	event, err := cc.eventService.GetEventById(c.Request().Context(), params.EventID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	user := ctx.User
	_, err = cc.eventRegistrationService.CreateEventRegistration(c.Request().Context(), user, event, &params)
	if err == nil {
		return c.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
	}

	var validationErrors validation.Errors
	ok := errors.As(err, &validationErrors)
	if ok {
		dto, err := cc.eventService.GetEventDetailsForEvent(c.Request().Context(), event, user)
		if err != nil {
			return err
		}

		c.Response().WriteHeader(http.StatusUnprocessableEntity)
		return eventregistrations.New(ctx, dto, &params, validationErrors).Render(c.Response())
	}

	return err
}
