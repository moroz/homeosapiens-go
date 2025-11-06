package handlers

import (
	"net/http"

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

	fs := http.FileServer(http.Dir("assets/dist"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	return r
}
