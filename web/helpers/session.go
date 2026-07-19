package helpers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/sessions"
)

func SaveSession(w http.ResponseWriter, store *sessions.Store, session sessions.Payload) error {
	cookie, err := store.EncodeSession(session)
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
		// Always consume the stored value, even when we reject it, so a poisoned
		// entry does not linger in the session.
		delete(ctx.Session, config.RedirectBackUrlSessionKey)
		if isSafeRedirectTarget(redirectTo) {
			return redirectTo
		}
	}
	return "/"
}

// isSafeRedirectTarget reports whether target is a local absolute path safe to
// redirect to after login. It rejects absolute URLs and protocol-relative URLs
// (e.g. "//evil.com", "/\\evil.com") to prevent open redirects.
func isSafeRedirectTarget(target string) bool {
	if !strings.HasPrefix(target, "/") || strings.HasPrefix(target, "//") || strings.HasPrefix(target, "/\\") {
		return false
	}
	u, err := url.Parse(target)
	if err != nil || u.Scheme != "" || u.Host != "" {
		return false
	}
	return true
}
