package handlers

import (
	"errors"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	userregistrations "github.com/moroz/homeosapiens-go/tmpl/user_registrations"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type userRegistrationController struct {
	srv *services.UserRegistrationService
}

func UserRegistrationController(db queries.DBTX) *userRegistrationController {
	return &userRegistrationController{
		srv: services.NewUserRegistrationService(db),
	}
}

func (cc *userRegistrationController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	return userregistrations.New(ctx, &types.RegisterUserParams{}, make(validation.Errors)).Render(c.Response())
}

func (cc *userRegistrationController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	var params types.RegisterUserParams
	if err := c.Bind(&params); err != nil {
		log.Print(err)
		return echo.ErrBadRequest
	}

	_, err := cc.srv.RegisterUser(c.Request().Context(), &params)
	if err, ok := errors.AsType[validation.Errors](err); ok {
		validationErrors := helpers.LocalizeValidationErrors(ctx.Localizer, err)
		return userregistrations.New(ctx, &params, validationErrors).Render(c.Response())
	}

	return nil
}
