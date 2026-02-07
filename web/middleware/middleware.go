package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/securecookie"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExtendContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Set("context", &types.CustomContext{})
		return next(c)
	}
}

func ResolveRequestLocale(bundle *goi18n.Bundle, store securecookie.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

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

func StoreRequestUrlInContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := helpers.GetRequestContext(c)
		ctx.RequestUrl = c.Request().URL
		ctx.RequestUrl.Host = c.Request().Host
		return next(c)
	}
}
