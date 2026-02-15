package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

func FetchCartFromSession(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			cartId, ok := ctx.Session["cart_id"].(string)
			if !ok {
				return next(c)
			}

			cart, err := queries.New(db).GetCart(c.Request().Context(), cartId)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			ctx.Cart = cart
			return next(c)
		}
	}
}
