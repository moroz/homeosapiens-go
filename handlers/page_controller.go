package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/pages"
)

type pageController struct {
	db           queries.DBTX
	eventService *services.EventService
}

func PageController(db queries.DBTX) *pageController {
	return &pageController{db, services.NewEventService(db)}
}

func (me *pageController) Index(c *echo.Context) error {
	ctx := c.Get("context").(*CustomContext)
	events, err := me.eventService.ListEvents(r.Context(), ctx.User)
	if err != nil {
		return err
	}

	if err := pages.Index(r.Context(), events).Render(w); err != nil {
		msg := fmt.Sprintf("Error rendering page: %s", err)
		log.Print(msg)
		http.Error(w, msg, 500)
	}
}
