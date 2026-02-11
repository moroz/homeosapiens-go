package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/admin"
	"github.com/moroz/homeosapiens-go/web/handlers"
	"github.com/moroz/homeosapiens-go/web/middleware"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Router(db queries.DBTX, bundle *i18n.Bundle, store securecookie.Store) http.Handler {
	r := echo.New()

	r.Pre(echomiddleware.MethodOverrideWithConfig(echomiddleware.MethodOverrideConfig{
		Getter: echomiddleware.MethodFromForm("_method"),
	}))
	r.Use(echomiddleware.RequestID())
	r.Use(echomiddleware.RequestLogger())
	r.Use(middleware.ExtendContext)
	r.Use(middleware.StoreRequestUrlInContext)
	r.Use(middleware.FetchSessionFromCookies(store, config.SessionCookieName))
	r.Use(middleware.FetchFlashMessages(store))
	r.Use(middleware.FetchUserFromSession(db))
	r.Use(middleware.ResolveTimezone)
	r.Use(middleware.ResolveRequestLocale(bundle, store))

	pages := handlers.PageController(db)
	r.GET("/", pages.Index)

	events := handlers.EventController(db)
	r.GET("/events/:slug", events.Show)

	eventRegistrations := handlers.EventRegistrationController(db)
	r.GET("/events/:slug/register", eventRegistrations.New)
	r.POST("/event_registrations", eventRegistrations.Create)

	sessions := handlers.SessionController(db, store)
	r.GET("/sign-in", sessions.New)
	r.POST("/sessions", sessions.Create)
	r.GET("/sign-out", sessions.Delete)

	userRegistrations := handlers.UserRegistrationController(db)
	r.GET("/sign-up", userRegistrations.New)

	videos := handlers.VideoController(db)
	r.GET("/videos", videos.Index)

	prefs := handlers.PreferencesController(store)
	r.POST("/api/v1/prefs/timezone", prefs.SaveTimezone)

	oauth2 := handlers.OAuth2Controller(store, db)
	r.GET("/oauth/google/redirect", oauth2.GoogleRedirect)
	r.GET("/oauth/google/callback", oauth2.GoogleCallback)

	auth := r.Group("")
	auth.Use(middleware.RequireAuthenticatedUser)

	profile := handlers.ProfileController(db)
	auth.GET("/profile", profile.Show)
	auth.PUT("/profile", profile.Update)

	ar := r.Group("/admin")
	ar.Use(middleware.RequireAdmin)

	adminEvents := admin.EventController(db)
	ar.GET("", adminEvents.Index)

	if config.IsProd {
		assets := r.Group("/assets")
		assets.Use(CacheControlMiddleware)
		assets.Static("", "assets/dist/assets")
	} else {
		r.Static("/assets", "assets/public/assets")

		email := handlers.EmailController()
		r.GET("/dev/email", email.Show)
	}

	return r
}

func CacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		path := c.Request().URL.Path

		// Cache versioned assets (containing hash in filename) for 1 year
		if strings.Contains(path, "-") && (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css")) {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Short cache for other assets
			c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		}

		return next(c)
	}
}
