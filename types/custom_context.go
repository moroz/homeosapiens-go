package types

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CustomContext struct {
	store       *session.Store
	Cart        *queries.GetCartRow
	CartId      *uuid.UUID
	Flash       Flash
	Language    string
	Localizer   *i18n.Localizer
	RequestUrl  *url.URL
	Session     session.Payload
	Timezone    *time.Location
	TimezoneSet bool
	User        *queries.User
}

func NewContext(store *session.Store) *CustomContext {
	return &CustomContext{store: store}
}

func (c *CustomContext) IsPolish() bool {
	return c.Language == "pl"
}

func (c *CustomContext) SaveSession(w http.ResponseWriter) error {
	cookie, err := c.store.EncodeSession(c.Session)
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

func (c *CustomContext) PutFlash(key, msg string) {
	c.Flash[key] = msg
	c.Session[config.FlashSessionKey] = c.Flash
}
