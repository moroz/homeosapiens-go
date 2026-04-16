package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/orders"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type orderController struct {
	cartService     *services.CartService
	eventService    *services.EventService
	orderService    *services.OrderService
	orderRepository *services.OrderRepository
}

func OrderController(db *pgxpool.Pool, service services.StripeService) *orderController {
	return &orderController{
		cartService:     services.NewCartService(db),
		eventService:    services.NewEventService(db),
		orderService:    services.NewOrderService(db, service),
		orderRepository: services.NewOrderRepository(db),
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

	return orders.New(ctx, cart, params, validation.Errors{}).Render(c.Response())
}

func (cc *orderController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	if ctx.CartId == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No cart ID in session")
	}

	cart, err := cc.cartService.GetCartItemsByCartId(c.Request().Context(), ctx.CartId)
	if err != nil {
		return err
	}
	if cart.IsEmpty() {
		return echo.NewHTTPError(http.StatusBadRequest, "Empty cart")
	}

	var params types.OrderParams
	if err := c.Bind(&params); err != nil {
		log.Print(err)
		return echo.ErrBadRequest
	}

	order, err := cc.orderService.CreateOrder(c.Request().Context(), *ctx.CartId, ctx.User, &params)
	if validationError, ok := errors.AsType[validation.Errors](err); ok {
		return orders.New(ctx, cart, &params, helpers.LocalizeValidationErrors(ctx.Localizer, validationError)).Render(c.Response())
	}

	return c.Redirect(http.StatusFound, order.CheckoutSession.URL)
}

func (cc *orderController) Success(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	sessionID := c.Request().URL.Query().Get("session_id")

	order, err := cc.orderRepository.GetOrderDetailsByCheckoutSessionID(c.Request().Context(), sessionID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Could not find an order with the checkout session ID %s.", sessionID))
	}

	if err != nil {
		return err
	}

	if order.PaidAt != nil {
		return orders.Success(ctx, order).Render(c.Response())
	} else {
		c.Response().Write([]byte("Order not paid"))
	}

	return nil
}
