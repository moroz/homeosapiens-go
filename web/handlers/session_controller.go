package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/sessions"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sessionController struct {
	db               queries.DBTX
	userService      *services.UserService
	userTokenService *services.UserTokenService
}

func SessionController(db queries.DBTX) *sessionController {
	return &sessionController{
		db:               db,
		userService:      services.NewUserService(db),
		userTokenService: services.NewUserTokenService(db),
	}
}

func (cc *sessionController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	redirectTo := c.QueryParams().Get("ref")
	if redirectTo != "" {
		ctx.Session[config.RedirectBackUrlSessionKey] = redirectTo
		_ = ctx.SaveSession(c.Response())
	}

	return sessions.New(ctx, "", "", false).Render(c.Response())
}

func (cc *sessionController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := cc.userService.AuthenticateUserByEmailPassword(c.Request().Context(), email, password)
	if err != nil {
		l := ctx.Localizer
		var msg string
		var msgIsHTML bool

		if errors.Is(err, services.ErrUnverifiedEmail) {
			resendURL := "/email-verifications/new?email=" + url.QueryEscape(email)
			msg = l.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    "sessions.new.unverified_email_html",
				TemplateData: map[string]string{"URL": resendURL},
			})
			msgIsHTML = true
		} else {
			msg = l.MustLocalizeMessage(&i18n.Message{ID: "sessions.new.invalid_email_password_combination"})
		}

		return sessions.New(ctx, email, msg, msgIsHTML).Render(c.Response())
	}

	return signUserIn(c, cc.db, user)
}

func (cc *sessionController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	token, ok := ctx.Session[config.AccessTokenSessionKey].([]byte)
	if ok && token != nil {
		if _, err := cc.userTokenService.DeleteToken(c.Request().Context(), token); err != nil {
			log.Printf("Error deleting user token: %s", err)
		}
	}
	delete(ctx.Session, config.AccessTokenSessionKey)
	if err := ctx.SaveSession(c.Response()); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
