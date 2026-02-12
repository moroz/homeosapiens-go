package helpers

import (
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
)

func SaveSession(w http.ResponseWriter, store securecookie.Store, session types.SessionData) error {
	bytes, err := session.Encode()
	if err != nil {
		return err
	}
	cookie, err := store.EncryptCookie(bytes)
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

func GetRedirectUrl(ctx *types.CustomContext) string {
	redirectTo, _ := ctx.Session[config.RedirectBackUrlSessionKey].(string)
	if redirectTo != "" {
		delete(ctx.Session, config.RedirectBackUrlSessionKey)
		return redirectTo
	}
	return "/"
}
