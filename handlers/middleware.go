package handlers

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocaleMiddleware(bundle *goi18n.Bundle, store securecookie.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(config.SessionContextName).(types.SessionData)

			langParam := r.FormValue("lang")
			header := r.Header.Get("Accept-Language")
			langFromSession, _ := session["lang"].(string)

			lang := i18n.ResolveLocale(langParam, langFromSession, header)

			if langParam != "" && langFromSession != langParam {
				storePreferredLangInSession(w, session, store, langParam)
			}

			localizer := goi18n.NewLocalizer(bundle, lang)
			ctx := context.WithValue(r.Context(), "localizer", localizer)
			ctx = context.WithValue(ctx, "lang", lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SaveSession(w http.ResponseWriter, store securecookie.Store, session types.SessionData) error {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(session)
	if err != nil {
		return err
	}
	cookie, err := store.EncryptCookie(buf.Bytes())
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    cookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}

func storePreferredLangInSession(w http.ResponseWriter, session types.SessionData, store securecookie.Store, newValue string) {
	session["lang"] = newValue
	_ = SaveSession(w, store, session)
}

func FetchSession(sessionStore securecookie.Store, cookieName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := decodeSessionFromRequest(sessionStore, cookieName, r)
			ctx := context.WithValue(r.Context(), config.SessionContextName, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FetchUserFromSession(db queries.DBTX) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(config.SessionContextName).(types.SessionData)
			var (
				user *queries.User
			)

			if token, ok := session["access_token"].([]byte); ok {
				user, _ = queries.New(db).GetUserByAccessToken(r.Context(), token)
			}

			ctx := context.WithValue(r.Context(), config.CurrentUserContextName, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func decodeSessionFromRequest(sessionStore securecookie.Store, cookieName string, r *http.Request) types.SessionData {
	result := make(types.SessionData)

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return result
	}

	binary, err := sessionStore.DecryptCookie(cookie.Value)
	if err != nil {
		return result
	}

	_ = gob.NewDecoder(bytes.NewBuffer(binary)).Decode(&result)

	return result
}
