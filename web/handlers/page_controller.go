package handlers

import (
	"sort"
	"time"

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
	allEvents, err := cc.eventService.ListPublishedEventsForUser(c.Request().Context(), ctx.User, ctx.CartId)
	if err != nil {
		return err
	}

	now := time.Now()
	var upcoming []*services.EventListDto
	for _, e := range allEvents {
		if e.EndsAt.After(now) {
			upcoming = append(upcoming, e)
		}
	}
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].StartsAt.Before(upcoming[j].StartsAt)
	})

	if len(upcoming) > 4 {
		upcoming = upcoming[:4]
	}

	return wrapRender(pages.Home(ctx, upcoming), c.Response())
}
