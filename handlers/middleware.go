package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/securecookie"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocaleMiddleware(bundle *goi18n.Bundle, store securecookie.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(config.SessionContextName).(*SessionData)

			langParam := r.FormValue("lang")
			header := r.Header.Get("Accept-Language")
			lang := i18n.ResolveLocale(langParam, session.Lang, header)

			if langParam != "" && session.Lang != langParam {
				storePreferredLangInSession(w, session, store, langParam)
			}

			localizer := goi18n.NewLocalizer(bundle, lang)
			ctx := context.WithValue(r.Context(), "localizer", localizer)
			ctx = context.WithValue(ctx, "lang", lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SaveSession(w http.ResponseWriter, store securecookie.Store, session *SessionData) error {
	asJson, err := json.Marshal(session)
	if err != nil {
		return err
	}
	cookie, err := store.EncryptCookie(asJson)
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

func storePreferredLangInSession(w http.ResponseWriter, session *SessionData, store securecookie.Store, newValue string) {
	session.Lang = newValue
	_ = SaveSession(w, store, session)
}

type SessionData struct {
	AccessToken []byte         `json:"access_token,omitempty"`
	Lang        string         `json:"lang,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
}

func FetchSession(sessionStore securecookie.Store, cookieName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := decodeSessionFromRequest(sessionStore, cookieName, r)
			ctx := context.WithValue(r.Context(), config.SessionContextName, &session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FetchUserFromSession(db queries.DBTX) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(config.SessionContextName).(*SessionData)
			var (
				user *queries.User
			)

			if session.AccessToken != nil {
				user, _ = queries.New(db).GetUserByAccessToken(r.Context(), session.AccessToken)
			}

			ctx := context.WithValue(r.Context(), config.CurrentUserContextName, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func decodeSessionFromRequest(sessionStore securecookie.Store, cookieName string, r *http.Request) (data SessionData) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}

	bytes, err := sessionStore.DecryptCookie(cookie.Value)
	if err != nil {
		return
	}

	_ = json.Unmarshal(bytes, &data)

	if data.Data == nil {
		data.Data = make(map[string]any)
	}
	return
}
