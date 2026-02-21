package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/carts"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type cartController struct {
	cartService  *services.CartService
	eventService *services.EventService
}

func CartController(db queries.DBTX) *cartController {
	return &cartController{
		cartService:  services.NewCartItemService(db),
		eventService: services.NewEventService(db),
	}
}

type createLineItemParams struct {
	EventId uuid.UUID `form:"event_id"`
}

func (cc *cartController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	var params createLineItemParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	event, err := cc.eventService.GetPaidEventById(c.Request().Context(), params.EventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cartItem, err := cc.cartService.AddEventToCart(c.Request().Context(), ctx.CartId(), params.EventId)
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

func (cc *cartController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	cart, err := cc.cartService.GetCartItemsByCartId(c.Request().Context(), ctx.CartId())
	if err != nil {
		return err
	}

	return carts.Show(ctx, cart).Render(c.Response())
}
