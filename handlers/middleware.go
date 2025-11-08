package handlers

import (
	"context"
	"net/http"

	"github.com/moroz/homeosapiens-go/i18n"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func AddI18NBundle(bundle *goi18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			langParam := r.FormValue("lang")
			header := r.Header.Get("Accept-Language")
			lang := i18n.ResolveLocale(langParam, header)
			localizer := goi18n.NewLocalizer(bundle, lang)
			ctx := context.WithValue(r.Context(), "localizer", localizer)
			ctx = context.WithValue(ctx, "lang", lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
