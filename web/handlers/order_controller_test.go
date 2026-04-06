package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/phone"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	count := func(ctx context.Context, table string) (int, error) {
		var val int
		err := db.QueryRow(ctx, "select count(*) from "+table).Scan(&val)
		return val, err
	}

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	mailer := mocks.NewMockMailer(t)
	mailer.EXPECT().Send(mock.Anything, mock.Anything).Return(nil)

	cs := mocks.GenerateCheckoutSession()

	stripeSrv := mocks.NewMockStripeService(t)
	stripeSrv.EXPECT().CreateCheckoutSession(mock.Anything, mock.Anything).Return(cs, nil)

	r := router.Router(db, store, stripeSrv, mailer)

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

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, 1, doc.Find(".cart-table").Length())

		var actions []string
		for _, f := range doc.Find("form[method=POST]").Nodes {
			for _, attr := range f.Attr {
				if attr.Key == "action" {
					actions = append(actions, attr.Val)
				}
			}
		}

		assert.Equal(t, 1, doc.Find("form[data-testid=checkout-form]").Length())
		assert.Zero(t, doc.Find("[data-testid=empty-message]").Length())
	})

	t.Run("cart view shows empty message when cart is empty", func(t *testing.T) {
		client, err := clientWithSession(store, origin, nil)
		require.NoError(t, err)

		resp, err := client.Get(srv.URL + "/cart")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		require.NoError(t, err)

		assert.Zero(t, doc.Find(".cart-table").Length())
		assert.Zero(t, doc.Find("form[data-testid=checkout-form]").Length())
		assert.NotZero(t, doc.Find("[data-testid=empty-message]").Length())
	})

	t.Run("placing order creates an order and a user", func(t *testing.T) {
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

		params := url.Values{
			"locale":                {"en"},
			"email":                 {"user@example.com"},
			"billing_address_line1": {"Example Street 42"},
			"billing_given_name":    {"John"},
			"billing_family_name":   {"Smith"},
			"billing_country":       {"DE"},
			"billing_city":          {"Berlin"},
			"billing_postal_code":   {"12345"},
			"billing_phone":         {phone.ExamplePhoneNumber("DE")},
		}
		body := bytes.NewBufferString(params.Encode())
		req, _ := http.NewRequest("POST", srv.URL+"/orders", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		countBefore, err := count(t.Context(), "orders")
		require.NoError(t, err)

		resp, err := client.Do(req)

		countAfter, err := count(t.Context(), "orders")
		require.NoError(t, err)
		assert.Equal(t, countBefore+1, countAfter)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)

		redirectedTo := resp.Header.Get("Location")
		assert.True(t, strings.HasPrefix(redirectedTo, "https://checkout.stripe.com/c/pay/"))
	})
}
