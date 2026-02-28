package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (queries.DBTX, error) {
	return pgxpool.New(ctx, config.MustGetenv("TEST_DATABASE_URL"))
}

func TestAddToCart(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)

	store, err := securecookie.NewStore(config.SessionKey)
	require.NoError(t, err)

	r := router.Router(db, store)
	require.NoError(t, err)

	srv := httptest.NewServer(r)
	defer srv.Close()

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
}
