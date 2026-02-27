package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type cartItemController struct {
	cartService  *services.CartService
	eventService *services.EventService
}

func CartItemController(db queries.DBTX) *cartItemController {
	return &cartItemController{
		cartService:  services.NewCartService(db),
		eventService: services.NewEventService(db),
	}
}

type lineItemParams struct {
	EventId uuid.UUID `form:"event_id"`
}

func (cc *cartItemController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	var params lineItemParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	event, err := cc.eventService.GetPaidEventById(c.Request().Context(), params.EventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cartItem, err := cc.cartService.AddEventToCart(c.Request().Context(), ctx.CartId, params.EventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ctx.CartId == nil {
		ctx.Session[config.CartIdSessionKey] = cartItem.CartID
		if err := ctx.SaveSession(c.Response()); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	redirectTo := c.Request().Referer()
	if redirectTo == "" {
		redirectTo = fmt.Sprintf("/events/%s", event.Slug)
	}

	return c.Redirect(http.StatusFound, redirectTo)
}

func (cc *cartItemController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	var params lineItemParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	if ctx.CartId != nil {
		found, err := cc.cartService.DeleteCartItem(c.Request().Context(), *ctx.CartId, params.EventId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if found {
			ctx.PutFlash("success", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
				ID: "cart_items.delete.success",
			}))
		}
	}

	return c.Redirect(http.StatusFound, "/cart")
}
