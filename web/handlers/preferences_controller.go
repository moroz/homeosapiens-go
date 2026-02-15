package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type preferencesController struct{}

func PreferencesController() *preferencesController {
	return &preferencesController{}
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

	c.Response().WriteHeader(http.StatusNoContent)
	return nil
}
