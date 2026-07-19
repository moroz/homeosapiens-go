package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	httpetag "github.com/go-http-utils/etag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/api"
	"github.com/moroz/homeosapiens-go/web/handlers"
	"github.com/moroz/homeosapiens-go/web/middleware"
	"github.com/moroz/homeosapiens-go/web/sessions"
)

type Groupie interface {
	Group(string, ...echo.MiddlewareFunc) *echo.Group
}

func Group(r Groupie, prefix string, cb func(r *echo.Group)) {
	group := r.Group(prefix)
	cb(group)
}

func Router(db *pgxpool.Pool, store *sessions.Store, stripeClient services.StripeService) *echo.Echo {
	r := echo.New()

	bundle, err := i18n.InitBundle()
	if err != nil {
		panic(err)
	}

	r.Pre(echomiddleware.MethodOverrideWithConfig(echomiddleware.MethodOverrideConfig{
		Getter: echomiddleware.MethodFromForm("_method"),
	}))
	// Reroute HEAD to GET for routing; the method is restored below before the handler runs.
	r.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Method == http.MethodHead {
				c.Set("_head", true)
				c.Request().Method = http.MethodGet
			}
			return next(c)
		}
	})
	r.Use(echomiddleware.RequestID())
	if !config.IsTest {
		r.Use(echomiddleware.RequestLogger())
	}
	r.Use(echomiddleware.Recover())

	if !config.IsProd {
		r.IPExtractor = echo.ExtractIPDirect()
		r.Static("/assets", "assets/public/assets")

		mountAdminSPAProxy(r)
	}

	r.Use(middleware.ExtendContext(store))
	r.Use(middleware.StoreRequestUrlInContext)
	// Restore HEAD method so the handler and logger both see the correct method,
	// and net/http suppresses the response body as required by the HTTP spec.
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Get("_head") != nil {
				c.Request().Method = http.MethodHead
			}
			return next(c)
		}
	})

	r.Use(middleware.FetchSessionFromCookies(store, config.SessionCookieName))
	r.Use(middleware.FetchFlashMessages(store))
	r.Use(middleware.FetchUserFromSession(db))
	r.Use(middleware.FetchCartFromSession(db))
	r.Use(middleware.ResolveTimezone)
	r.Use(middleware.ResolveRequestLocale(bundle))
	r.Use(middleware.StoreLocalePreference(db))

	pages := handlers.PageController(db)
	r.GET("/", pages.Index)

	events := handlers.EventController(db)
	r.GET("/events", events.Index)
	r.GET("/events/:slug", events.Show)

	prefs := handlers.PreferencesController(db)
	r.POST("/api/v1/prefs/timezone", prefs.SaveTimezone)

	oauth2 := handlers.OAuth2Controller(db)
	r.GET("/oauth/google/redirect", oauth2.GoogleRedirect)
	r.GET("/oauth/google/callback", oauth2.GoogleCallback)

	orders := handlers.OrderController(db, stripeClient)
	r.GET("/cart", orders.New)
	r.POST("/orders", orders.Create)
	r.GET("/orders/success", orders.Success)

	stripe := handlers.StripeWebhookController(db, stripeClient)
	r.POST("/webhooks/stripe", echo.WrapHandler(http.HandlerFunc(stripe.StripeWebhook)))

	cartItems := handlers.CartItemController(db)
	r.POST("/cart_items", cartItems.Create)
	r.DELETE("/cart_items", cartItems.Delete)

	sessions := handlers.SessionController(db)
	r.DELETE("/sign-out", sessions.Delete)

	videos := handlers.VideoController(db)
	r.GET("/videos/:id/thumbnail/:locale", videos.Thumbnail, echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return httpetag.Handler(h, false)
	}))
	r.GET("/watch", videos.Youtube)

	Group(r, "", func(r *echo.Group) {
		r.Use(middleware.RedirectToHomeIfAuthenticated)

		r.GET("/sign-in", sessions.New)
		r.POST("/sessions", sessions.Create)

		userRegistrations := handlers.UserRegistrationController(db)
		r.GET("/sign-up", userRegistrations.New)
		r.POST("/sign-up", userRegistrations.Create)
		r.GET("/user-registrations/success", userRegistrations.Success)

		userPasswords := handlers.UserPasswordResetController(db)
		r.GET("/reset-password", userPasswords.New)
		r.POST("/reset-password", userPasswords.Create)

		Group(r, "", func(r *echo.Group) {
			r.Use(userPasswords.VerifyToken)

			r.GET("/reset-password/:token", userPasswords.Edit)
			r.PUT("/reset-password/:token", userPasswords.Update)
		})

		emailVerifications := handlers.EmailVerificationController(db)
		r.GET("/email-verifications/new", emailVerifications.New)
		r.POST("/email-verifications", emailVerifications.Create)
		r.GET("/verify-email", emailVerifications.Verify)
	})

	Group(r, "", func(r *echo.Group) {
		r.Use(middleware.RequireAuthenticatedUser)

		profile := handlers.ProfileController(db)
		r.GET("/profile", profile.Show)
		r.PUT("/profile", profile.Update)

		eventRegistrations := handlers.EventRegistrationController(db)
		r.GET("/events/:event_id/register", eventRegistrations.Create)
		r.POST("/event_registrations/:event_id", eventRegistrations.Create)
		r.DELETE("/event_registrations/:event_id", eventRegistrations.Delete)

		videos := handlers.VideoController(db)
		r.GET("/videos/:group_slug/:video_slug", videos.Show)
		r.GET("/videos/:group_slug", videos.Index)
		r.GET("/videos", videos.Index)
	})

	// JSON admin API (OpenAPI/oapi-codegen). Routes are generated from
	// web/api/openapi.yaml and served under /api/admin. Kept off the /admin
	// prefix so the whole /admin/* namespace belongs to the SPA with no
	// proxy exceptions.
	Group(r, "/api/admin", func(r *echo.Group) {
		r.Use(middleware.RequireAdmin)

		apiServer := api.NewServer(db)
		r.Any("/*", echo.WrapHandler(apiServer.Handler("/api/admin")))
	})

	// Admin SPA (React Router, library mode). In dev, reverse-proxy the whole
	// /admin/* namespace to the Vite dev server so HMR works. In prod, Caddy
	// serves the built dist/ directly and never reaches Go.
	if !config.IsProd {
	}

	if !config.IsProd {
		email := handlers.EmailController(db)
		r.GET("/dev/email/order", email.OrderConfirmation)
		r.GET("/dev/email/payment", email.PaymentConfirmation)
		r.GET("/dev/email/email_verification", email.UserEmailVerification)
		r.GET("/dev/email/password_reset", email.PasswordReset)
		r.GET("/dev/email/event_registration", email.EventRegistrationConfirmation)
	}

	return r
}

// mountAdminSPAProxy reverse-proxies the whole /admin/* namespace to the Vite
// dev server (default http://localhost:5173, override with VITE_DEV_SERVER).
// Vite is configured with base "/admin/", so the path is forwarded unchanged.
// Dev only — in prod Caddy serves the built assets and Go never sees /admin.
func mountAdminSPAProxy(r *echo.Echo) {
	target, err := url.Parse(config.ViteDevServer)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	handler := echo.WrapHandler(proxy)

	r.Any("/admin/*", handler, middleware.RequireAdmin)

	r.GET("/admin", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/admin/")
	})
}
