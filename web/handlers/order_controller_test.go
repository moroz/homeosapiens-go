package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.MustGetenv("TEST_DATABASE_URL"))
}

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

	t.Run("add to cart", func(t *testing.T) {
		jar, err := cookiejar.New(nil)
		require.NoError(t, err)

		params := url.Values{
			"event_id": {"019c5c9a-c5a4-7518-8317-65ae90516726"},
		}
		body := bytes.NewBufferString(params.Encode())

		client := &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		req, _ := http.NewRequest("POST", srv.URL+"/cart_items", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)

		origin, err := url.Parse(srv.URL)
		require.NoError(t, err)

		var cookie *http.Cookie
		for _, c := range jar.Cookies(origin) {
			if c.Name == config.SessionCookieName {
				cookie = c
			}
		}
		assert.NotNil(t, cookie)

		sessionPayload, err := store.DecodeSession(cookie)
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

	})
}
