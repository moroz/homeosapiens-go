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
	r.Use(FetchUserFromSession(db))
	r.Use(LocaleMiddleware(bundle, store))

	pages := PageController(db)
	r.Get("/", pages.Index)

	events := EventController(db)
	r.Get("/events/{slug}", events.Show)

	sessions := SessionController(db, store)
	r.Get("/sign-in", sessions.New)
	r.Post("/sessions", sessions.Create)

	r.Handle("/assets/*", MultiDirFileServer("assets/dist", "assets/static"))

	return r
}

// MultiDirFileServer handles requests for static assets backed by multiple directories.
// TODO: Make this development-only
func MultiDirFileServer(dirs ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for _, d := range dirs {
			stripped, _ := strings.CutPrefix(path, "/assets/")
			f, err := http.Dir(d).Open(stripped)
			if err == nil {
				f.Close()
				http.StripPrefix("/assets/", http.FileServer(http.Dir(d))).ServeHTTP(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})
}
