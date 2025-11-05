package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()

	pages := PageController()
	r.Get("/", pages.Index)
	return r
}
