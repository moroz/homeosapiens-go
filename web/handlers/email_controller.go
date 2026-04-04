package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type emailController struct {
	db           queries.DBTX
	orderService *services.OrderService
}

func EmailController(db queries.DBTX) *emailController {
	return &emailController{
		orderService: services.NewOrderService(db, nil),
		db:           db,
	}
}

func (cc *emailController) getLastOrder(c *echo.Context) (*types.OrderDTO, error) {
	lastID, err := queries.New(cc.db).GetLastOrderID(c.Request().Context())
	if err != nil {
		return nil, err
	}

	return cc.orderService.GetOrderDetails(c.Request().Context(), lastID)
}

func (cc *emailController) OrderConfirmation(c *echo.Context) error {
	order, err := cc.getLastOrder(c)
	if err != nil {
		return err
	}

	ctx := helpers.GetRequestContext(c)
	props := email.OrderConfirmationEmailProps{
		LayoutProps: &email.LayoutProps{
			LogoURL:   config.PublicUrl + "/assets/logo.png",
			Localizer: ctx.Localizer,
			Language:  ctx.Language,
		},
		Order: order,
	}

	return email.OrderConfirmationTemplate.Execute(c.Response(), props)
}

func (cc *emailController) PaymentConfirmation(c *echo.Context) error {
	order, err := cc.getLastOrder(c)
	if err != nil {
		return err
	}

	ctx := helpers.GetRequestContext(c)
	props := email.PaymentConfirmationEmailProps{
		LayoutProps: &email.LayoutProps{
			LogoURL:   config.PublicUrl + "/assets/logo.png",
			Localizer: ctx.Localizer,
			Language:  ctx.Language,
		},
		Order: order,
	}

	return email.PaymentConfirmationTemplate.Execute(c.Response(), props)
}
