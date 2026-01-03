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
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type oauth2Controller struct {
	sessionStore securecookie.Store
	config       *oauth2.Config
	*services.UserService
	*services.UserTokenService
}

func OAuth2Controller(store securecookie.Store, db queries.DBTX) *oauth2Controller {
	return &oauth2Controller{
		store,
		&oauth2.Config{
			ClientID:     config.GoogleClientId,
			ClientSecret: config.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  config.PublicUrl + "/oauth/google/callback",
			Scopes:       []string{"email", "profile"},
		},
		services.NewUserService(db),
		services.NewUserTokenService(db),
	}
}

const OAuth2SessionKey = "auth_state"
const RedirectBackUrlSessionKey = "redirect_back"

func (c *oauth2Controller) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	if c.config.ClientID == "" {
		log.Printf("Google Client ID is not set")
		http.Error(w, "Client ID is not set", http.StatusInternalServerError)
		return
	}

	redirectTo := r.URL.Query().Get("ref")

	var state = make([]byte, 4)
	_, _ = rand.Read(state)
	session := r.Context().Value(config.SessionContextName).(types.SessionData)
	session[OAuth2SessionKey] = hex.EncodeToString(state)
	if redirectTo != "" {
		session[RedirectBackUrlSessionKey] = redirectTo
	}
	if err := SaveSession(w, c.sessionStore, session); err != nil {
		log.Printf("Error persisting session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := c.config.AuthCodeURL(hex.EncodeToString(state), oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func decodeIDTokenClaims(token string) (*types.GoogleIDTokenClaims, error) {
	segs := strings.Split(token, ".")
	bytes, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		return nil, err
	}
	var claims types.GoogleIDTokenClaims
	if err := json.Unmarshal(bytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (c *oauth2Controller) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionContextName).(types.SessionData)
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
		log.Printf("Google token exchange returned error: %s", err)
		http.Error(w, "Failed to fetch access token", 500)
		return
	}

	idToken, _ := token.Extra("id_token").(string)

	validator, _ := idtoken.NewValidator(r.Context())
	_, err = validator.Validate(r.Context(), idToken, c.config.ClientID)
	if err != nil {
		log.Printf("ID token verification failed: %s", err)
		http.Error(w, fmt.Sprintf("ID token verification failed: %s", err), 500)
		return
	}

	claims, err := decodeIDTokenClaims(idToken)
	if err != nil {
		log.Printf("Failed to decode ID token: %s", err)
		http.Error(w, fmt.Sprintf("Failed to decode ID token: %s", err), 500)
		return
	}

	user, err := c.UserService.FindOrCreateUserFromClaims(r.Context(), claims)
	if err != nil {
		log.Printf("Failed to create user from claims: %s", err)
		http.Error(w, fmt.Sprintf("Failed to create user from claims: %s", err), 500)
		return
	}

	userToken, err := c.UserTokenService.IssueAccessTokenForUser(r.Context(), user, 24*time.Hour)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	redirectBackUrl, ok := session[RedirectBackUrlSessionKey].(string)
	if !ok {
		redirectBackUrl = "/"
	}

	session["access_token"] = userToken.Token
	delete(session, OAuth2SessionKey)
	delete(session, RedirectBackUrlSessionKey)
	if err := SaveSession(w, c.sessionStore, session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, redirectBackUrl, http.StatusFound)
}
