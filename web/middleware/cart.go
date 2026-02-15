package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

func FetchCartFromSession(db queries.DBTX) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			// Check if there is a cart ID stored in the session, or if the user is signed in
			// (the user may have a default cart)
			var cartId *uuid.UUID
			if id, ok := ctx.Session["cart_id"].(uuid.UUID); ok {
				cartId = &id
			}

			var userId *uuid.UUID
			if ctx.User != nil {
				userId = &ctx.User.ID
			}

			// If there is no cart ID or user in the session, there can't possibly be a cart
			if cartId == nil && userId == nil {
				return next(c)
			}

			cart, err := queries.New(db).GetCart(c.Request().Context(), &queries.GetCartParams{
				CartID:  cartId,
				OwnerID: userId,
			})
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			ctx.Cart = cart
			return next(c)
		}
	}
}
