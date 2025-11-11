package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Router(db queries.DBTX, bundle *i18n.Bundle) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(AddI18NBundle(bundle))

	pages := PageController(db)
	r.Get("/", pages.Index)

	events := EventController(db)
	r.Get("/events/{slug}", events.Show)

	r.Handle("/assets/*", MultiDirFileServer("assets/dist", "assets/static"))

	return r
}

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
