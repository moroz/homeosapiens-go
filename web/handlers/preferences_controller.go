package handlers

import (
	"log"
	"net/http"
	"time"

	sqlcrypter "github.com/bincyber/go-sqlcrypter"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type preferencesController struct {
	db queries.DBTX
}

func PreferencesController(db queries.DBTX) *preferencesController {
	return &preferencesController{db: db}
}

func (cc *preferencesController) SaveTimezone(c *echo.Context) error {
	tzParam := c.QueryParam("tz")
	if _, err := time.LoadLocation(tzParam); err != nil || tzParam == "" {
		return echo.NewHTTPError(400, "Invalid timezone")
	}

	ctx := helpers.GetRequestContext(c)
	ctx.Session["tz"] = tzParam
	if err := ctx.SaveSession(c.Response()); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		return err
	}

	if ctx.User != nil {
		existing := ctx.User.PreferredTimezone
		if existing == nil || existing.String() != tzParam {
			enc := sqlcrypter.NewEncryptedBytes(tzParam)
			if err := queries.New(cc.db).UpdateUserPreferredTimezone(c.Request().Context(), &queries.UpdateUserPreferredTimezoneParams{
				PreferredTimezone: &enc,
				ID:                ctx.User.ID,
			}); err != nil {
				log.Printf("Error persisting timezone preference: %s", err)
			}
		}
	}

	c.Response().WriteHeader(http.StatusNoContent)
	return nil
}
