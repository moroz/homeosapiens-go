package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
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

	eventId := c.Param("event_id")
	if eventId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required parameter: event_id")
	}

	event, err := cc.eventService.GetEventById(c.Request().Context(), eventId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	_, err = cc.eventRegistrationService.CreateEventRegistration(c.Request().Context(), ctx.User, event)
	if err == nil {
		return c.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
	}

	return err
}

func (cc *eventRegistrationController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	eventId := c.Param("event_id")
	if eventId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required parameter: event_id")
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
