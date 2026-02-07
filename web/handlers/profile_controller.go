package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/profile"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type profileController struct {
	userService *services.UserService
}

func ProfileController(db queries.DBTX) *profileController {
	return &profileController{
		userService: services.NewUserService(db),
	}
}

func (cc *profileController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	return profile.Show(ctx).Render(c.Response())
}
