package admin

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/admin/users"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type userController struct {
	userService *services.UserService
}

func UserController(db queries.DBTX) *userController {
	return &userController{
		userService: services.NewUserService(db),
	}
}

func (cc *userController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	list, err := cc.userService.ListUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return users.Index(ctx, list).Render(c.Response())
}
