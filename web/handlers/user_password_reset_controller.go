package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/user_password_resets"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type userPasswordResetController struct {
	srv services.UserPasswordResetService
}

func UserPasswordResetController(db queries.DBTX) *userPasswordResetController {
	return &userPasswordResetController{*services.NewUserPasswordResetService(db)}
}

// New displays a form to the user which can be used to request a password reset token with an email address.
func (cc *userPasswordResetController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	return user_password_resets.New(ctx, "", "").Render(c.Response())
}

// Create is the submission handler for the form displayed by New.
func (cc *userPasswordResetController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	_ = ctx

	panic("unimplemented")
}

// Edit displays the form for the user to update their password. The URL contains a password reset token. If the token provided in the URL is invalid or expired, access is denied.
func (cc *userPasswordResetController) Edit(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	_ = ctx

	token := c.Param("token")
	_ = token

	panic("unimplemented")
}

func (cc *userPasswordResetController) Update(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	_ = ctx

	token := c.Param("token")
	_ = token

	panic("unimplemented")
}
