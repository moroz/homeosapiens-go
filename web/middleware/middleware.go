package middleware

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/securecookie"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func storePreferredLangInSession(w http.ResponseWriter, session types.SessionData, store securecookie.Store, newValue string) {
	session["lang"] = newValue
	_ = helpers.SaveSession(w, store, session)
}

func decodeSessionFromRequest(sessionStore securecookie.Store, cookieName string, r *http.Request) types.SessionData {
	result := make(types.SessionData)

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return result
	}

	binary, err := sessionStore.DecryptCookie(cookie.Value)
	if err != nil {
		return result
	}

	_ = gob.NewDecoder(bytes.NewBuffer(binary)).Decode(&result)

	return result
}

func ExtendContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Set("context", &types.CustomContext{})
		return next(c)
	}
}

func FetchSessionFromCookies(sessionStore securecookie.Store, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*types.CustomContext)
			ctx.Session = decodeSessionFromRequest(sessionStore, cookieName, c.Request())
			return next(c)
		}
	}
}

func ResolveRequestLocale(bundle *goi18n.Bundle, store securecookie.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*types.CustomContext)

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

func FetchUserFromSession(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Get("context").(*types.CustomContext)

			if token, ok := ctx.Session["access_token"].([]byte); ok {
				if u, err := queries.New(db).GetUserByAccessToken(c.Request().Context(), token); err == nil {
					ctx.User = u
				}
			}

			return next(c)
		}
	}
}

func ResolveTimezone(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Get("context").(*types.CustomContext)
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

func StoreRequestUrlInContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Get("context").(*types.CustomContext)
		ctx.RequestUrl = c.Request().URL
		ctx.RequestUrl.Host = c.Request().Host
		return next(c)
	}
}
