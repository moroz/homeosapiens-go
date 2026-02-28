package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/orders"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type orderController struct {
	cartService *services.CartService
}

func OrderController(db queries.DBTX) *orderController {
	return &orderController{
		cartService: services.NewCartService(db),
	}
}

func (cc *orderController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	cart, err := cc.cartService.GetCartItemsByCartId(c.Request().Context(), ctx.CartId)
	if err != nil {
		return err
	}

	if cart.IsEmpty() {
		ctx.PutFlash("error", "Your cart is currently empty.")
		if err := ctx.SaveSession(c.Response()); err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/cart")
	}

	return orders.New(ctx).Render(c.Response())
}
