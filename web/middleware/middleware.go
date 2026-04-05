package middleware

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"github.com/moroz/homeosapiens-go/web/session"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExtendContext(store *session.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("context", types.NewContext(store))
			return next(c)
		}
	}
}

func ResolveRequestLocale(bundle *goi18n.Bundle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			langParam := c.FormValue("lang")
			header := c.Request().Header.Get("Accept-Language")
			langFromSession, _ := ctx.Session["lang"].(string)

			lang := i18n.ResolveLocale(langParam, langFromSession, header)

			if langParam != "" && langFromSession != langParam {
				ctx.Session["lang"] = langParam
				if err := ctx.SaveSession(c.Response()); err != nil {
					return err
				}
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

func StoreLocalePreference(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)
			if ctx.User == nil {
				return next(c)
			}

			if ctx.User.PreferredLocale != queries.Locale(ctx.Language) {
				go queries.New(db).UpdateUserPreferredLocale(context.Background(), &queries.UpdateUserPreferredLocaleParams{
					PreferredLocale: queries.Locale(ctx.Language),
					ID:              ctx.User.ID,
				})
			}

			return next(c)
		}
	}
}
