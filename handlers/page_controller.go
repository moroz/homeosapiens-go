package handlers

import (
	"fmt"
	"log"
	"net/http"

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

func (c *pageController) Index(w http.ResponseWriter, r *http.Request) {
	events, err := c.eventService.ListEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := pages.Index(r.Context(), events).Render(w); err != nil {
		msg := fmt.Sprintf("Error rendering page: %s", err)
		log.Print(msg)
		http.Error(w, msg, 500)
	}
}
