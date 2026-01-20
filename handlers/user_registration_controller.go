package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	userregistrations "github.com/moroz/homeosapiens-go/tmpl/user_registrations"
	"github.com/moroz/homeosapiens-go/types"
)

type userRegistrationController struct {
	*services.UserService
}

func UserRegistrationController(db queries.DBTX) *userRegistrationController {
	return &userRegistrationController{
		services.NewUserService(db),
	}
}

func (c *userRegistrationController) New(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)
	return userregistrations.New(ctx).Render(r.Response())
}
