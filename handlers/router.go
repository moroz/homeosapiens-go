package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Router(db queries.DBTX, bundle *i18n.Bundle, store securecookie.Store) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(FetchSession(store, config.SessionCookieName))
	r.Use(FetchPreferredTimezone)
	r.Use(StoreRequestUrlInContext)
	r.Use(FetchUserFromSession(db))
	r.Use(LocaleMiddleware(bundle, store))

	oauth2 := OAuth2Controller(store, db)
	r.Get("/oauth/google/redirect", oauth2.GoogleRedirect)
	r.Get("/oauth/google/callback", oauth2.GoogleCallback)

	return r
}

// CacheControlMiddleware adds cache-control headers based on file patterns
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Cache versioned assets (containing hash in filename) for 1 year
		if strings.Contains(path, "-") && (strings.HasSuffix(path, ".js") ||
			strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".woff2")) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Short cache for other assets
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}

		next.ServeHTTP(w, r)
	})
}
