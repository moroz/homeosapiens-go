package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/dashboard"
	"github.com/moroz/homeosapiens-go/types"
)

type dashboardController struct {
	db queries.DBTX
}

func DashboardController(db queries.DBTX) *dashboardController {
	return &dashboardController{db}
}

func (me *dashboardController) Index(c *echo.Context) error {
	ctx := c.Get("context").(*types.CustomContext)
	return dashboard.Index(ctx).Render(c.Response())
}
