package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsNew(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)
	r := router.Router(db, store, mocks.NewMockStripeService(t))

	server := httptest.NewServer(r)
	defer server.Close()

	origin, err := url.Parse(server.URL)
	require.NoError(t, err)

	_, err = db.Exec(t.Context(), "truncate users cascade")
	user, err := mocks.UniqueUser(db, t.Context(), func(params *types.SeedUserParams) {
		params.Email = "session-test@example.com"
		params.Password = "foobar2000"
		params.EmailConfirmed = true
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.EmailConfirmedAt)

	inputs := []struct {
		email, password string
		valid           bool
	}{
		{"session-test@example.com", "foobar2000", true},
		{"SeSsIoN-TEST@example.com", "foobar2000", true},
		{"SESSION-TEST@EXAMPLE.COM", "foobar2000", true},
		{"SESSION-TEST@EXAPLE.COM", "foobar2000", false},
		{"session", "foobar2000", false},
		{"", "foobar2000", false},
	}

	for _, input := range inputs {
		client, err := clientWithSession(store, origin, nil)

		body := bytes.NewBufferString(url.Values{
			"email":    {input.email},
			"password": {input.password},
		}.Encode())

		req, err := http.NewRequest("POST", server.URL+"/sessions", body)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		require.NoError(t, err)

		clientSession, err := getClientSession(client.Jar, store, origin)
		if input.valid {
			require.NoError(t, err)
			assert.NotEmpty(t, clientSession)
			assert.NotEmpty(t, resp.Header.Get("Location"))
			assert.Equal(t, http.StatusFound, resp.StatusCode)
		} else {
			require.Error(t, err)
			assert.Empty(t, clientSession)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp.Body.Close()
	}
}
