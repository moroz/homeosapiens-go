package helpers

import (
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
)

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
