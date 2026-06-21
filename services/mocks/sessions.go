package mocks

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/sessions"
)

func ClientWithSession(store *sessions.Store, origin *url.URL, payload sessions.Payload) (*http.Client, error) {
	var cookie string

	if payload != nil {
		sessionCookie, err := store.EncodeSession(payload)
		if err != nil {
			return nil, err
		}

		cookie = sessionCookie
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if cookie != "" {
		jar.SetCookies(origin, []*http.Cookie{
			{
				Name:     config.SessionCookieName,
				Value:    cookie,
				Secure:   true,
				HttpOnly: true,
			},
		})
	}

	return &http.Client{Jar: jar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}, nil
}

func GetClientSession(jar http.CookieJar, store *sessions.Store, origin *url.URL) (sessions.Payload, error) {
	var cookie *http.Cookie
	for _, c := range jar.Cookies(origin) {
		if c.Name == config.SessionCookieName {
			cookie = c
		}
	}
	if cookie == nil {
		return nil, errors.New("session cookie not found")
	}

	return store.DecodeSession(cookie)
}

type ClientWithUserInput struct {
	Store   *sessions.Store
	Server  *httptest.Server
	User    *queries.User
	DB      queries.DBTX
	Context context.Context
}

func ClientWithUser(props *ClientWithUserInput) (*http.Client, error) {
	origin, err := url.Parse(props.Server.URL)
	if err != nil {
		return nil, err
	}

	var payload sessions.Payload
	if props.User != nil {
		payload, err = UserSession(props.DB, props.Context, props.User)
		if err != nil {
			return nil, err
		}
	}

	return ClientWithSession(props.Store, origin, payload)
}
