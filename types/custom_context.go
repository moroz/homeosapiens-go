package types

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/url"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CustomContext struct {
	store       securecookie.Store
	User        *queries.User
	Session     SessionData
	Localizer   *i18n.Localizer
	Language    string
	Timezone    *time.Location
	TimezoneSet bool
	RequestUrl  *url.URL
	Flash       Flash
}

func NewContext(store securecookie.Store) *CustomContext {
	return &CustomContext{store: store}
}

func (c *CustomContext) IsPolish() bool {
	return c.Language == "pl"
}

func (c *CustomContext) SaveSession(w http.ResponseWriter) error {
	b, err := c.Session.Encode()
	if err != nil {
		return err
	}

	cookie, err := c.store.EncryptCookie(b)
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

func (c *CustomContext) DecodeSession(cookie *http.Cookie) SessionData {
	result := make(SessionData)

	if cookie == nil {
		return result
	}

	binary, err := c.store.DecryptCookie(cookie.Value)
	if err != nil {
		return result
	}
	_ = gob.NewDecoder(bytes.NewBuffer(binary)).Decode(&result)
	return result
}

func (s SessionData) Encode() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
