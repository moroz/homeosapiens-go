package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/middleware"
	"github.com/moroz/securecookie"
)

type preferencesController struct {
	sessionStore securecookie.Store
}

func PreferencesController(sessionStore securecookie.Store) *preferencesController {
	return &preferencesController{
		sessionStore,
	}
}

func (c *preferencesController) SaveTimezone(r *echo.Context) error {
	tzParam := r.QueryParam("tz")
	if _, err := time.LoadLocation(tzParam); err != nil || tzParam == "" {
		return echo.NewHTTPError(400, "Invalid timezone")
	}
	ctx := r.Get("context").(*types.CustomContext)
	ctx.Session["tz"] = tzParam
	if err := middleware.SaveSession(r.Response(), c.sessionStore, ctx.Session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		return err
	}

	r.Response().WriteHeader(http.StatusNoContent)
	return nil
}
