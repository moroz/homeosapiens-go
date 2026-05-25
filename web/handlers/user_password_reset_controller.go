package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/user_password_resets"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type userPasswordResetController struct {
	userTokenService     *services.UserTokenService
	passwordResetService *services.UserPasswordResetService
}

func UserPasswordResetController(db queries.DBTX) *userPasswordResetController {
	return &userPasswordResetController{
		userTokenService:     services.NewUserTokenService(db),
		passwordResetService: services.NewUserPasswordResetService(db),
	}
}

// New displays a form to the user which can be used to request a password reset token with an email address.
func (cc *userPasswordResetController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	return user_password_resets.New(ctx, "", "").Render(c.Response())
}

// Create is the submission handler for the form displayed by New.
func (cc *userPasswordResetController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	l := ctx.Localizer

	var params types.SendPasswordResetParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest
	}

	sent, err := cc.passwordResetService.MaybeIssuePasswordResetTokenForUser(c.Request().Context(), params.Email)
	if errors.Is(err, services.ErrRateLimited) {
		ctx.PutFlash("error", l.MustLocalizeMessage(&i18n.Message{
			ID: "user_password_resets.create.rate_limited",
		}))
	} else if sent || errors.Is(err, services.ErrUserNonExistent) {
		ctx.PutHTMLFlash("success", l.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "user_password_resets.create.success_flash",
			TemplateData: map[string]string{
				"Email": params.Email,
			},
		}))
	} else if err != nil {
		log.Printf("Error in UserPasswordResetController.Create: %v", err)
		return echo.ErrInternalServerError
	}

	if err := ctx.SaveSession(c.Response()); err != nil {
		log.Printf("Error saving session: %v", err)
	}

	return c.Redirect(http.StatusFound, "/sign-in")
}

// Edit displays the form for the user to update their password. The URL contains a password reset token. If the token provided in the URL is invalid or expired, access is denied.
func (cc *userPasswordResetController) Edit(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	_ = ctx

	token := c.Param("token")
	_, err := cc.userTokenService.FindUserByUserTokenParam(c.Request().Context(), token, "password_reset")
	if err != nil {
		ctx.PutFlash("error", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
			ID: "user_password_resets.edit.token_invalid_or_expired",
		}))

		if err := ctx.SaveSession(c.Response()); err != nil {
			log.Printf("Error saving session: %v", err)
			return echo.ErrInternalServerError
		}

		return c.Redirect(http.StatusFound, "/")
	}

	panic("unimplemented")
}

func (cc *userPasswordResetController) Update(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	_ = ctx

	token := c.Param("token")
	_ = token

	panic("unimplemented")
}
