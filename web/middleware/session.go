package middleware

import (
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/securecookie"
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

func FetchSessionFromCookies(sessionStore securecookie.Store, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)
			ctx.Session = decodeSessionFromRequest(sessionStore, cookieName, c.Request())
			return next(c)
		}
	}
}

func FetchFlashMessages(sessionStore securecookie.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			if flash, ok := ctx.Session["_flash"].([]types.FlashMessage); ok {
				delete(ctx.Session, "_flash")
				_ = helpers.SaveSession(c.Response(), sessionStore, ctx.Session)
				ctx.Flash = flash
			}

			return next(c)
		}
	}
}
