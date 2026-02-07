package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/sessions"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sessionController struct {
	*services.UserService
	*services.UserTokenService
	sessionStore securecookie.Store
}

func SessionController(db queries.DBTX, sessionStore securecookie.Store) *sessionController {
	return &sessionController{
		services.NewUserService(db),
		services.NewUserTokenService(db),
		sessionStore,
	}
}

func (cc *sessionController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	redirectTo := c.QueryParams().Get("ref")
	if redirectTo != "" {
		ctx.Session[config.RedirectBackUrlSessionKey] = redirectTo
		_ = helpers.SaveSession(c.Response(), cc.sessionStore, ctx.Session)
	}

	return sessions.New(ctx, "", "").Render(c.Response())
}

func (cc *sessionController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := cc.UserService.AuthenticateUserByEmailPassword(c.Request().Context(), email, password)
	if err != nil {
		l := ctx.Localizer
		msg := l.MustLocalizeMessage(&i18n.Message{
			ID: "sessions.new.invalid_email_password_combination",
		})

		return sessions.New(ctx, email, msg).Render(c.Response())
	}

	token, err := cc.UserTokenService.IssueAccessTokenForUser(c.Request().Context(), user, 24*time.Hour)
	if err != nil {
		return err
	}

	ctx.Session["access_token"] = token.Token
	redirectTo := helpers.GetRedirectUrl(ctx)

	if err := helpers.SaveSession(c.Response(), cc.sessionStore, ctx.Session); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, redirectTo)
}

func (cc *sessionController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	token, ok := ctx.Session["access_token"].([]byte)
	if ok && token != nil {
		if _, err := cc.UserTokenService.DeleteToken(c.Request().Context(), token); err != nil {
			log.Printf("Error deleting user token: %s", err)
		}
	}
	delete(ctx.Session, "access_token")
	if err := helpers.SaveSession(c.Response(), cc.sessionStore, ctx.Session); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
