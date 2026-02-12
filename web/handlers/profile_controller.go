package handlers

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/profile"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (cc *profileController) Update(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	var params types.UpdateProfileRequest
	if err := c.Bind(&params); err != nil {
		log.Print(err)
		return echo.ErrBadRequest
	}

	user, err := cc.userService.UpdateUserProfile(c.Request().Context(), ctx.User, &params)
	if err != nil {
		log.Print(err)
		return echo.ErrInternalServerError
	}

	ctx.PutFlash("success", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
		ID: "profile.messages.success",
	}))

	ctx.User = user
	return profile.Show(ctx).Render(c.Response())
}
