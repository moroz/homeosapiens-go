package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth2Controller struct {
	store securecookie.Store
}

func OAuth2Controller(store securecookie.Store) *oauth2Controller {
	return &oauth2Controller{store}
}

func (c *oauth2Controller) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	var state = make([]byte, 4)
	_, _ = rand.Read(state)
	session := r.Context().Value(config.SessionContextName).(*SessionData)
	session.Data["auth_state"] = hex.EncodeToString(state)
	if err := SaveSession(w, c.store, session); err != nil {
		log.Printf("Error persisting session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conf := &oauth2.Config{
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000/oauth/google/callback",
		Scopes:       []string{"email", "profile"},
	}

	url := conf.AuthCodeURL(hex.EncodeToString(state), oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}
