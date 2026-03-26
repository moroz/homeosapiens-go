package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/orders"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type orderController struct {
	cartService  *services.CartService
	eventService *services.EventService
}

func OrderController(db queries.DBTX) *orderController {
	return &orderController{
		cartService:  services.NewCartService(db),
		eventService: services.NewEventService(db),
	}
}

func (cc *orderController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	cart, err := cc.cartService.GetCartItemsByCartId(c.Request().Context(), ctx.CartId)
	if err != nil {
		return err
	}

	params := &types.OrderParams{
		BillingCountry: "PL",
	}

	if ctx.User != nil && ctx.User.Country != nil {
		params.BillingCountry = *ctx.User.Country
	} else if tzGuess := countries.GuessRegionByTimezone(ctx.Timezone.String()); tzGuess.Found && countries.IsEUMemberState(tzGuess.IsoCode) {
		params.BillingCountry = tzGuess.IsoCode
	}

	if ctx.User != nil {
		params.BillingGivenName = ctx.User.GivenName.String()
		params.BillingFamilyName = ctx.User.FamilyName.String()
		params.Email = ctx.User.Email.String()
	}

	return orders.New(ctx, cart, params).Render(c.Response())
}
