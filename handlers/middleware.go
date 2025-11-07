package handlers

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func AddI18NBundle(bundle *i18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.FormValue("lang")
			accept := r.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(bundle, lang, accept)
			ctx := context.WithValue(r.Context(), "localizer", localizer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
