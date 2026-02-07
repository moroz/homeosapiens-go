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

		return next(c)
	}
}
