package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (cc *eventRegistrationController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	eventId, err := uuid.Parse(c.Param("event_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid required parameter: event_id")
	}

	event, err := cc.eventService.GetRegisterableEventById(c.Request().Context(), eventId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	registered, err := cc.eventRegistrationService.CreateEventRegistration(c.Request().Context(), ctx.User, event)
	if registered {
		ctx.PutFlash("success", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
			ID: "event_registrations.create.success",
		}))
	}
	_ = ctx.SaveSession(c.Response())

	if err == nil {
		return c.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
	}

	return err
}

func (cc *eventRegistrationController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	eventId, err := uuid.Parse(c.Param("event_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid required parameter: event_id")
	}

	event, err := cc.eventService.GetEventById(c.Request().Context(), eventId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	_, err = cc.eventRegistrationService.DeleteEventRegistration(c.Request().Context(), ctx.User, event)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
}
