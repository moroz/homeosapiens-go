package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/carts"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type cartController struct {
	cartService  *services.CartService
	eventService *services.EventService
}

func CartController(db queries.DBTX) *cartController {
	return &cartController{
		cartService:  services.NewCartService(db),
		eventService: services.NewEventService(db),
	}
}

func (cc *cartController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	cart, err := cc.cartService.GetCartItemsByCartId(c.Request().Context(), ctx.CartId)
	if err != nil {
		return err
	}

	params := &types.OrderParams{
		BillingCountry: "PL",
	}

	if tzGuess := countries.GuessRegionByTimezone(ctx.Timezone.String()); tzGuess.Found && countries.IsEUMemberState(tzGuess.IsoCode) {
		params.BillingCountry = tzGuess.IsoCode
	}

	return carts.Show(ctx, cart, params).Render(c.Response())
}
