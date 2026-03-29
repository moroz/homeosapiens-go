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
)

type emailController struct {
	db           queries.DBTX
	orderService *services.OrderService
}

func EmailController(db queries.DBTX) *emailController {
	return &emailController{
		orderService: services.NewOrderService(db, services.MockStripeService()),
		db:           db,
	}
}

func (c *emailController) Show(r *echo.Context) error {
	lastID, err := queries.New(c.db).GetLastOrderID(r.Request().Context())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "No Order found in the DB")
	}

	order, err := c.orderService.GetOrderDetails(r.Request().Context(), lastID)
	if err != nil {
		return err
	}

	props := email.OrderConfirmationEmailProps{
		LayoutProps: &email.LayoutProps{
			LogoURL: config.PublicUrl + "/assets/logo.png",
		},
		Order: order,
	}

	return email.OrderConfirmationTemplate.Execute(r.Response(), props)
}
