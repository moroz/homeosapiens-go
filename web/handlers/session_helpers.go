package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

func signUserIn(c *echo.Context, db queries.DBTX, user *queries.User) error {
	ctx := helpers.GetRequestContext(c)

	token, err := services.NewUserTokenService(db).IssueAccessTokenForUser(c.Request().Context(), user, config.AccessTokenValidity)
	if err != nil {
		return err
	}

	_ = services.NewUserService(db).SetUserLastLogin(c.Request().Context(), c.RealIP(), user.ID)
	ctx.Session[config.AccessTokenSessionKey] = token.Token
	ctx.Session[config.LanguageSessionKey] = string(user.PreferredLocale)
	delete(ctx.Session, config.OAuth2SessionKey)
	delete(ctx.Session, config.RedirectBackUrlSessionKey)
	redirectTo := helpers.GetRedirectUrl(ctx)

	if err := ctx.SaveSession(c.Response()); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, redirectTo)
}
