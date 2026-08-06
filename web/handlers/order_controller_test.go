package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/internal/phone"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.MustGetenv("TEST_DATABASE_URL"))
}

func TestCartFlow(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	paidEvent, err := mocks.PaidEvent(db, t.Context())
	require.NoError(t, err)
	paidEventId := paidEvent.ID

	count := func(ctx context.Context, table string) (int, error) {
		var val int
		err := db.QueryRow(ctx, "select count(*) from "+table).Scan(&val)
		return val, err
	}

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	cs := mocks.CheckoutSession()

	stripeSrv := mocks.NewMockStripeService(t)
	stripeSrv.EXPECT().CreateCheckoutSession(mock.Anything, mock.Anything).Return(cs, nil)

	r := router.Router(db, store, stripeSrv)

	srv := httptest.NewServer(r)
	defer srv.Close()

	origin, err := url.Parse(srv.URL)
	require.NoError(t, err)

	t.Run("add to cart", func(t *testing.T) {
		params := url.Values{
			"event_id": {paidEventId.String()},
		}
		body := bytes.NewBufferString(params.Encode())

		client, err := mocks.ClientWithSession(store, origin, nil)
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", srv.URL+"/cart_items", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)

		sessionPayload, err := mocks.GetClientSession(client.Jar, store, origin)
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
		_, err = db.Exec(t.Context(), "insert into cart_line_items (cart_id, product_id) select $1, e.product_id from events e where e.id = $2", cartId, paidEventId)
		require.NoError(t, err)

		client, err := mocks.ClientWithSession(store, origin, sessions.Payload{
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
		client, err := mocks.ClientWithSession(store, origin, nil)
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
		_, err = db.Exec(t.Context(), "insert into cart_line_items (cart_id, product_id) select $1, e.product_id from events e where e.id = $2", cartId, paidEventId)
		require.NoError(t, err)

		client, err := mocks.ClientWithSession(store, origin, sessions.Payload{
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

		_, err = db.Exec(t.Context(), "truncate river_job")
		require.NoError(t, err)

		resp, err := client.Do(req)

		countAfter, err := count(t.Context(), "orders")
		require.NoError(t, err)
		assert.Equal(t, countBefore+1, countAfter)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)

		redirectedTo := resp.Header.Get("Location")
		assert.True(t, strings.HasPrefix(redirectedTo, "https://checkout.stripe.com/c/pay/"))

		rivertest.RequireInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendOrderEmailArgs{}, nil)
	})
}

// Attendance to a past event is worthless, so an ended event can neither be put
// in the cart nor paid for if it ended while sitting there.
func TestCartRejectsEndedEvents(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	srv := httptest.NewServer(router.Router(db, store, mocks.NewMockStripeService(t)))
	defer srv.Close()

	origin, err := url.Parse(srv.URL)
	require.NoError(t, err)

	endedEvent, err := mocks.PaidEvent(db, ctx, func(p *queries.UpsertEventParams) {
		p.StartsAt = time.Now().UTC().Add(-26 * time.Hour)
		p.EndsAt = time.Now().UTC().Add(-24 * time.Hour)
	})
	require.NoError(t, err)

	t.Run("POST /cart_items refuses an event that has ended", func(t *testing.T) {
		client, err := mocks.ClientWithSession(store, origin, nil)
		require.NoError(t, err)

		body := bytes.NewBufferString(url.Values{"event_id": {endedEvent.ID.String()}}.Encode())
		req, _ := http.NewRequest("POST", srv.URL+"/cart_items", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("POST /orders refuses a cart holding an event that has ended", func(t *testing.T) {
		cartId := uuid.Must(uuid.NewV7())
		_, err = db.Exec(ctx, "insert into cart_line_items (cart_id, product_id) select $1, e.product_id from events e where e.id = $2", cartId, endedEvent.ID)
		require.NoError(t, err)

		client, err := mocks.ClientWithSession(store, origin, sessions.Payload{
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
		req, _ := http.NewRequest("POST", srv.URL+"/orders", bytes.NewBufferString(params.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		countBefore, err := countRows(ctx, db, "orders")
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.Equal(t, "/cart", resp.Header.Get("Location"))

		countAfter, err := countRows(ctx, db, "orders")
		require.NoError(t, err)
		assert.Equal(t, countBefore, countAfter)
	})
}

func countRows(ctx context.Context, db *pgxpool.Pool, table string) (int, error) {
	var val int
	err := db.QueryRow(ctx, "select count(*) from "+table).Scan(&val)
	return val, err
}
