package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/email"
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

func (cc *emailController) Show(c *echo.Context) error {
	lastID, err := queries.New(cc.db).GetLastOrderID(c.Request().Context())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "No Order found in the DB")
	}

	order, err := cc.orderService.GetOrderDetails(c.Request().Context(), lastID)
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
