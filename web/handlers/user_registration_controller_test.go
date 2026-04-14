package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyEmail(t *testing.T) {
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

	service := services.NewUserTokenService(db)

	t.Run("verifies user and signs them in with valid token", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context())
		require.NoError(t, err)
		require.Nil(t, user.EmailConfirmedAt)

		token, err := service.IssueEmailVerificationTokenForUser(t.Context(), user, config.EmailVerificationTokenValidity)
		require.NoError(t, err)
		require.NotNil(t, token)

		assert.Nil(t, user.EmailConfirmedAt)

		client, err := clientWithSession(store, origin, nil)
		req, err := http.NewRequest("GET", server.URL+token.VerifyEmailPath(), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// account activated
		user, err = queries.New(db).GetUserByID(t.Context(), user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.EmailConfirmedAt)

		sessionPayload, err := getClientSession(client.Jar, store, origin)
		assert.NotNil(t, sessionPayload)
		assert.NoError(t, err)

		// user signed in
		accessToken := sessionPayload[config.AccessTokenSessionKey]
		assert.NotEmpty(t, accessToken)
	})

	t.Run("returns an unprocessable entity with expired token", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context())
		require.NoError(t, err)
		require.Nil(t, user.EmailConfirmedAt)

		token, err := service.IssueEmailVerificationTokenForUser(t.Context(), user, config.EmailVerificationTokenValidity)
		require.NoError(t, err)
		require.NotNil(t, token)
		assert.Nil(t, user.EmailConfirmedAt)

		_, err = db.Exec(t.Context(), "update user_tokens set inserted_at = now() - interval '5 days', valid_until = now() - interval '1 day' where id = $1", token.UserToken.ID)
		require.NoError(t, err)

		client, err := clientWithSession(store, origin, nil)
		req, err := http.NewRequest("GET", server.URL+token.VerifyEmailPath(), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		user, err = queries.New(db).GetUserByID(t.Context(), user.ID)
		require.NoError(t, err)
		assert.Nil(t, user.EmailConfirmedAt)
	})
}
