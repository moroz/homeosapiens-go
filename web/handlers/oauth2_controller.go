package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/helpers"
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

func (cc *oauth2Controller) GoogleRedirect(c *echo.Context) error {
	if cc.config.ClientID == "" {
		log.Printf("Google Client ID is not set")
		return echo.NewHTTPError(500, "Client ID is not set")
	}

	redirectTo := c.QueryParam("ref")

	ctx := helpers.GetRequestContext(c)

	var state = make([]byte, 4)
	_, _ = rand.Read(state)
	ctx.Session[OAuth2SessionKey] = hex.EncodeToString(state)
	if redirectTo != "" {
		ctx.Session[RedirectBackUrlSessionKey] = redirectTo
	}
	if err := helpers.SaveSession(c.Response(), cc.sessionStore, ctx.Session); err != nil {
		log.Printf("Error persisting session: %s", err)
		return err
	}

	url := cc.config.AuthCodeURL(hex.EncodeToString(state), oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
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

func (cc *oauth2Controller) GoogleCallback(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	state, _ := ctx.Session[OAuth2SessionKey].(string)
	stateParam := c.QueryParam("state")

	if state != stateParam {
		log.Printf("Invalid OAuth2 state param in callback")
		return echo.NewHTTPError(400, "Invalid OAuth2 state param")
	}

	code := c.QueryParam("code")
	token, err := cc.config.Exchange(c.Request().Context(), code)
	if err != nil {
		log.Printf("Google token exchange returned error: %s", err)

		return echo.NewHTTPError(500, "Failed to fetch access token")
	}

	idToken, _ := token.Extra("id_token").(string)

	validator, _ := idtoken.NewValidator(c.Request().Context())
	_, err = validator.Validate(c.Request().Context(), idToken, cc.config.ClientID)
	if err != nil {
		log.Printf("ID token verification failed: %s", err)
		return err
	}

	claims, err := decodeIDTokenClaims(idToken)
	if err != nil {
		log.Printf("Failed to decode ID token: %s", err)
		return err
	}

	user, err := cc.UserService.FindOrCreateUserFromClaims(c.Request().Context(), claims)
	if err != nil {
		log.Printf("Failed to create user from claims: %s", err)
		return err
	}

	userToken, err := cc.UserTokenService.IssueAccessTokenForUser(c.Request().Context(), user, 24*time.Hour)
	if err != nil {
		log.Printf("Error issuing access token: %s", err)
		return err
	}

	redirectBackUrl, ok := ctx.Session[RedirectBackUrlSessionKey].(string)
	if !ok {
		redirectBackUrl = "/"
	}

	ctx.Session["access_token"] = userToken.Token
	delete(ctx.Session, OAuth2SessionKey)
	delete(ctx.Session, RedirectBackUrlSessionKey)
	if err := helpers.SaveSession(c.Response(), cc.sessionStore, ctx.Session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		return err
	}

	return c.Redirect(http.StatusFound, redirectBackUrl)
}
