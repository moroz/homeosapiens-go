package middleware

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

func ResolveTimezone(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := helpers.GetRequestContext(c)
		ctx.Timezone, _ = time.LoadLocation("Europe/Warsaw")

		if tzFromSession, ok := ctx.Session["tz"].(string); ok && tzFromSession != "" {
			if loaded, err := time.LoadLocation(tzFromSession); err == nil {
				ctx.Timezone = loaded
				ctx.TimezoneSet = true
			}
		}

		// For logged-in users with no session timezone yet, fall back to the
		// value persisted in the database from a previous session.
		if !ctx.TimezoneSet && ctx.User != nil && ctx.User.PreferredTimezone != nil {
			if loaded, err := time.LoadLocation(ctx.User.PreferredTimezone.String()); err == nil {
				ctx.Timezone = loaded
				ctx.TimezoneSet = true
			}
		}

		return next(c)
	}
}
