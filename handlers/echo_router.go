package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func EchoRouter(db queries.DBTX, bundle *i18n.Bundle, store securecookie.Store) http.Handler {
	r := echo.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(ExtendContext)
	r.Use(FetchSessionEcho(store, config.SessionCookieName))
	r.Use(FetchUserFromSessionEcho(db))
	r.Use(FetchPreferredTimezoneEcho)
	r.Use(LocaleMiddlewareEcho(bundle, store))

	pages := PageController(db)
	r.GET("/", pages.Index)

	return r
}
