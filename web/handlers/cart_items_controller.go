package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type cartItemsController struct {
	cartService  *services.CartService
	eventService *services.EventService
}

func CartItemsController(db queries.DBTX) *cartItemsController {
	return &cartItemsController{
		cartService:  services.NewCartItemService(db),
		eventService: services.NewEventService(db),
	}
}

type createLineItemParams struct {
	EventId uuid.UUID `form:"event_id"`
}

func (cc *cartItemsController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	var params createLineItemParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	event, err := cc.eventService.GetPaidEventById(c.Request().Context(), params.EventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cartItem, err := cc.cartService.AddEventToCart(c.Request().Context(), ctx.CartId(), ctx.User, params.EventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ctx.CartId() == nil {
		ctx.Session[config.CartIdSessionKey] = cartItem.CartID
		if err := ctx.SaveSession(c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/events/%s", event.Slug))
}
