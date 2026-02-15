package admin

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/admin/events"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type eventController struct {
	eventService *services.EventService
}

func EventController(db queries.DBTX) *eventController {
	return &eventController{
		eventService: services.NewEventService(db),
	}
}

func (cc *eventController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	list, err := cc.eventService.ListEvents(c.Request().Context(), ctx.User)
	if err != nil {
		return err
	}

	return events.Index(ctx, list).Render(c.Response())
}

func (cc *eventController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid required param: id")
	}

	event, err := cc.eventService.GetEventDetailsById(c.Request().Context(), id, ctx.User)
	if err != nil {
		return err
	}

	return events.Show(ctx, event).Render(c.Response())
}
