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
	r.Use(StoreRequestUrlInContextEcho)
	r.Use(FetchSessionEcho(store, config.SessionCookieName))
	r.Use(FetchUserFromSessionEcho(db))
	r.Use(FetchPreferredTimezoneEcho)
	r.Use(LocaleMiddlewareEcho(bundle, store))

	pages := PageController(db)
	r.GET("/", pages.Index)

	events := EventController(db)
	r.GET("/events/:slug", events.Show)

	eventRegistrations := EventRegistrationController(db)
	r.GET("/events/:slug/register", eventRegistrations.New)
	r.POST("/event_registrations", eventRegistrations.Create)

	sessions := SessionController(db, store)
	r.GET("/sign-in", sessions.New)
	r.POST("/sessions", sessions.Create)
	r.GET("/sign-out", sessions.Delete)

	userRegistrations := UserRegistrationController(db)
	r.GET("/sign-up", userRegistrations.New)

	videos := VideoController(db)
	r.GET("/videos", videos.Index)

	prefs := PreferencesController(store)
	r.POST("/api/v1/prefs/timezone", prefs.SaveTimezone)

	if config.IsProd {
		fileServer := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/dist/assets")))
		r.GET("/assets/*", echo.WrapHandler(fileServer), echo.WrapMiddleware(CacheControlMiddleware))
	} else {
		email := EmailController()
		r.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/public/assets")))))
		r.GET("/dev/email", email.Show)
	}

	return r
}
