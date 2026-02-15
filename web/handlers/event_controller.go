package handlers

import (
	"database/sql"
	"errors"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type eventController struct {
	db           queries.DBTX
	eventService *services.EventService
}

func EventController(db queries.DBTX) *eventController {
	return &eventController{db, services.NewEventService(db)}
}

func (cc *eventController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	slug := c.Param("slug")
	user := ctx.User
	event, err := cc.eventService.GetEventDetailsBySlug(c.Request().Context(), slug, user, ctx.CartId())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	if err != nil {
		log.Printf("Error fetching event by slug %s: %s", slug, err)
		return err
	}

	return events.Show(ctx, event).Render(c.Response())
}
