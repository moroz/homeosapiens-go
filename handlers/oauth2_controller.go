package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth2Controller struct {
	store  securecookie.Store
	config *oauth2.Config
}

func OAuth2Controller(store securecookie.Store) *oauth2Controller {
	return &oauth2Controller{
		store,
		&oauth2.Config{
			ClientID:     config.GoogleClientId,
			ClientSecret: config.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:3000/oauth/google/callback",
			Scopes:       []string{"email", "profile"},
		},
	}
}

const OAuth2SessionKey = "auth_state"

func (c *oauth2Controller) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	var state = make([]byte, 4)
	_, _ = rand.Read(state)
	session := r.Context().Value(config.SessionContextName).(SessionData)
	session[OAuth2SessionKey] = hex.EncodeToString(state)
	if err := SaveSession(w, c.store, session); err != nil {
		log.Printf("Error persisting session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := c.config.AuthCodeURL(hex.EncodeToString(state), oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (c *oauth2Controller) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionContextName).(SessionData)
	state, _ := session[OAuth2SessionKey].(string)
	stateParam := r.URL.Query().Get("state")

	if state != stateParam {
		log.Printf("Invalid OAuth2 state param in callback")
		http.Error(w, "Invalid state", 400)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := c.config.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Google token exchange returned error: %s", token)
		http.Error(w, "Failed to fetch access token", 500)
		return
	}

	fmt.Fprintf(w, "%+v", token)
}
