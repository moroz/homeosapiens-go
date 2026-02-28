package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.MustGetenv("TEST_DATABASE_URL"))
}

func noFollow(r *http.Request, v []*http.Request) error {
	return http.ErrUseLastResponse
}

func getClientSession(jar http.CookieJar, store *session.Store, origin *url.URL) (session.Payload, error) {
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

func clientWithSession(store *session.Store, origin *url.URL, payload session.Payload) (*http.Client, error) {
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

	return &http.Client{Jar: jar, CheckRedirect: noFollow}, nil
}

var PaidEventId = uuid.MustParse("019c5c9a-c5a4-7518-8317-65ae90516726")

func TestCartFlow(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	r := router.Router(db, store)
	require.NoError(t, err)

	srv := httptest.NewServer(r)
	defer srv.Close()

	origin, err := url.Parse(srv.URL)
	require.NoError(t, err)

	t.Run("add to cart", func(t *testing.T) {
		params := url.Values{
			"event_id": {PaidEventId.String()},
		}
		body := bytes.NewBufferString(params.Encode())

		client, err := clientWithSession(store, origin, nil)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", srv.URL+"/cart_items", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)

		sessionPayload, err := getClientSession(client.Jar, store, origin)
		assert.NotNil(t, sessionPayload)
		assert.NoError(t, err)

		cartId, ok := sessionPayload[config.CartIdSessionKey].(uuid.UUID)
		assert.True(t, ok)
		assert.NotEqual(t, uuid.UUID{}, cartId)

		cart, err := queries.New(db).GetCart(t.Context(), cartId)
		assert.NoError(t, err)
		assert.True(t, cart.ProductTotal.Equal(decimal.NewFromInt(560)))
		assert.Equal(t, int64(1), cart.ItemCount)
		assert.Equal(t, cartId, cart.CartID)
	})

	t.Run("cart view", func(t *testing.T) {
		cartId := uuid.Must(uuid.NewV7())
		item, err := queries.New(db).InsertCartLineItem(t.Context(), &queries.InsertCartLineItemParams{
			CartID:  cartId,
			EventID: PaidEventId,
		})
		assert.NoError(t, err)
		assert.NotNil(t, item)

		client, err := clientWithSession(store, origin, session.Payload{
			config.CartIdSessionKey: cartId,
		})
		require.NoError(t, err)

		resp, err := client.Get(srv.URL + "/cart")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		html, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NotEmpty(t, html)
	})

	t.Run("checkout redirects to /cart if cart is empty", func(t *testing.T) {
		client, err := clientWithSession(store, origin, nil)
		require.NoError(t, err)

		resp, err := client.Get(srv.URL + "/checkout")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)

		assert.Equal(t, "/cart", resp.Header.Get("Location"))

		sessionPayload, err := getClientSession(client.Jar, store, origin)
		assert.NoError(t, err)
		assert.NotNil(t, sessionPayload)

		flash, ok := sessionPayload[config.FlashSessionKey].(types.Flash)
		assert.True(t, ok)
		assert.NotNil(t, flash)

		assert.NotEmpty(t, flash["error"])
	})
}
