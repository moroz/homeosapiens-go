package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/pages"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type pageController struct {
	db           queries.DBTX
	eventService *services.EventService
}

func PageController(db queries.DBTX) *pageController {
	return &pageController{db, services.NewEventService(db)}
}

func (cc *pageController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	events, err := cc.eventService.ListEvents(c.Request().Context(), ctx.User, ctx.CartId)
	if err != nil {
		return err
	}

	return pages.Index(ctx, events).Render(c.Response())
}
