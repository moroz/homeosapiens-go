package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/homeosapiens-go/web/session"
)

func FetchSessionFromCookies(store *session.Store, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)
			cookie, _ := c.Cookie(cookieName)
			payload, _ := store.DecodeSession(cookie)
			ctx.Session = payload
			return next(c)
		}
	}
}

func FetchFlashMessages(sessionStore *session.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)
			ctx.Flash = make(types.Flash)

			if flash, ok := ctx.Session[config.FlashSessionKey].(types.Flash); ok {
				delete(ctx.Session, config.FlashSessionKey)
				_ = helpers.SaveSession(c.Response(), sessionStore, ctx.Session)
				ctx.Flash = flash
			}

			return next(c)
		}
	}
}
