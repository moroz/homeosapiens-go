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

func Group(r *echo.Echo, prefix string, cb func(r *echo.Group)) {
	group := r.Group(prefix)
	cb(group)
}

func Router(db queries.DBTX, bundle *i18n.Bundle, store securecookie.Store) http.Handler {
	r := echo.New()

	r.Pre(echomiddleware.MethodOverrideWithConfig(echomiddleware.MethodOverrideConfig{
		Getter: echomiddleware.MethodFromForm("_method"),
	}))
	r.Use(echomiddleware.RequestID())
	r.Use(echomiddleware.RequestLogger())

	if config.IsProd {
		r.IPExtractor = echo.ExtractIPFromXFFHeader()
		r.Static("/assets", "assets/dist/assets", CacheControlMiddleware)
	} else {
		r.IPExtractor = echo.ExtractIPDirect()
		r.Static("/assets", "assets/public/assets")
	}

	r.Use(middleware.ExtendContext(store))
	r.Use(middleware.StoreRequestUrlInContext)

	r.Use(middleware.FetchSessionFromCookies(config.SessionCookieName))
	r.Use(middleware.FetchFlashMessages(store))
	r.Use(middleware.FetchUserFromSession(db))
	r.Use(middleware.ResolveTimezone)
	r.Use(middleware.ResolveRequestLocale(bundle))

	pages := handlers.PageController(db)
	r.GET("/", pages.Index)

	events := handlers.EventController(db)
	r.GET("/events/:slug", events.Show)

	eventRegistrations := handlers.EventRegistrationController(db)
	r.GET("/events/:slug/register", eventRegistrations.New)
	r.POST("/event_registrations", eventRegistrations.Create)

	sessions := handlers.SessionController(db)
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

	Group(r, "", func(r *echo.Group) {
		r.Use(middleware.RequireAuthenticatedUser)

		profile := handlers.ProfileController(db)
		r.GET("/profile", profile.Show)
		r.PUT("/profile", profile.Update)
	})

	Group(r, "/admin", func(r *echo.Group) {
		r.Use(middleware.RequireAdmin)

		events := admin.EventController(db)
		r.GET("", events.Index)
		r.GET("/events/:id", events.Show)

		users := admin.UserController(db)
		r.GET("/users", users.Index)
	})

	if !config.IsProd {
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
