package handlers

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type cartItemsController struct {
	cartService *services.CartService
}

type createLineItemParams struct {
	EventId pgtype.UUID `form:"event_id"`
}

func (cc *cartItemsController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

}
