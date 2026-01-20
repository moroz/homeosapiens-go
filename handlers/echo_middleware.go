package handlers

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/securecookie"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExtendContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Set("context", &CustomContext{})
		return next(c)
	}
}

func FetchSessionEcho(sessionStore securecookie.Store, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*CustomContext)
			ctx.Session = decodeSessionFromRequest(sessionStore, cookieName, c.Request())
			return next(c)
		}
	}
}

func LocaleMiddlewareEcho(bundle *goi18n.Bundle, store securecookie.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*CustomContext)

			langParam := c.FormValue("lang")
			header := c.Request().Header.Get("Accept-Language")
			langFromSession, _ := ctx.Session["lang"].(string)

			lang := i18n.ResolveLocale(langParam, langFromSession, header)

			if langParam != "" && langFromSession != langParam {
				storePreferredLangInSession(c.Response(), ctx.Session, store, langParam)
			}

			ctx.Localizer = goi18n.NewLocalizer(bundle, lang)
			ctx.Language = lang
			return next(c)
		}
	}
}

func FetchUserFromSessionEcho(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*CustomContext)

			if token, ok := ctx.Session["access_token"].([]byte); ok {
				if u, err := queries.New(db).GetUserByAccessToken(r.Context(), token); err == nil {
					ctx.User = u
				}
			}

			return next(c)
		}
	}
}

func FetchPreferredTimezoneEcho(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Get("context").(*CustomContext)
		ctx.Timezone, _ = time.LoadLocation("Europe/Warsaw")

		if tzFromSession, ok := ctx.Session["tz"].(string); ok && tzFromSession != "" {
			if loaded, err := time.LoadLocation(tzFromSession); err == nil {
				ctx.Timezone = loaded
				ctx.TimezoneSet = true
			}
		}

		return next(c)
	}
}
