package middleware

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

func FetchUserFromSession(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			if token, ok := ctx.Session[config.AccessTokenSessionKey].([]byte); ok {
				if u, err := queries.New(db).GetUserByAccessToken(c.Request().Context(), &queries.GetUserByAccessTokenParams{
					Token:   token,
					Context: config.UserTokenContextAccess,
				}); err == nil {
					ctx.User = u
				}
			}

			return next(c)
		}
	}
}

func FetchUserFromBearerToken(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			if ctx.User != nil {
				return next(c)
			}

			header := c.Request().Header.Get("Authorization")
			if token, found := strings.CutPrefix(header, "Bearer "); found {
				if u, err := queries.New(db).GetUserByAccessToken(c.Request().Context(), &queries.GetUserByAccessTokenParams{
					Token:   crypto.HashUserToken([]byte(token)),
					Context: config.UserTokenContextApiKey,
				}); err == nil {
					ctx.User = u
					ctx.IsApiToken = true
				}
			}

			return next(c)
		}
	}
}

func RequireAuthenticatedUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := helpers.GetRequestContext(c)

		if ctx.User != nil {
			return next(c)
		}

		qs := url.Values{
			"ref": {ctx.RequestUrl.String()},
		}
		ctx.PutFlash("danger", "You need to sign in to access this page.")
		if err := ctx.SaveSession(c.Response()); err != nil {
			log.Print(err)
		}

		return c.Redirect(http.StatusSeeOther, "/sign-in?"+qs.Encode())
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := helpers.GetRequestContext(c)

		if ctx.User != nil && ctx.User.UserRole == queries.UserRoleAdministrator {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}
}

func RedirectToHomeIfAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := helpers.GetRequestContext(c)

		if ctx.User == nil {
			return next(c)
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}
}
