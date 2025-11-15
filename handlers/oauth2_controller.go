package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
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

type IDTokenClaims struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Avatar     string `json:"picture"`
}

func decodeIDTokenClaims(token string) (*IDTokenClaims, error) {
	segs := strings.Split(token, ".")
	bytes, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		return nil, err
	}
	var claims IDTokenClaims
	if err := json.Unmarshal(bytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
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

	id_token, _ := token.Extra("id_token").(string)

	validator, _ := idtoken.NewValidator(r.Context())
	_, err = validator.Validate(r.Context(), id_token, c.config.ClientID)
	if err != nil {
		log.Printf("ID token verification failed: %s", err)
		http.Error(w, fmt.Sprintf("ID token verification failed: %s", err), 500)
		return
	}

	claims, err := decodeIDTokenClaims(id_token)
	if err != nil {
		log.Printf("Failed to decode ID token: %s", err)
		http.Error(w, fmt.Sprintf("Failed to decode ID token: %s", err), 500)
		return
	}

	pretty, _ := json.MarshalIndent(claims, "", "\t")
	w.Write(pretty)
}
