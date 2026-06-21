package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRegistrationController(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)

	_, err = db.Exec(ctx, "truncate event_registrations")
	require.NoError(t, err)

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	r := router.Router(db, store, mocks.NewMockStripeService(t))

	server := httptest.NewServer(r)
	defer server.Close()

	countRegistrations := func() (int, error) {
		var count int
		err := db.QueryRow(ctx, "select count(*) from event_registrations").Scan(&count)
		return count, err
	}

	t.Run("when signed in", func(t *testing.T) {
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		client, err := mocks.ClientWithUser(&mocks.ClientWithUserInput{
			Store:   store,
			Server:  server,
			User:    user,
			DB:      db,
			Context: ctx,
		})
		require.NoError(t, err)

		t.Run("POST /event_registrations/:event_id signs up for free event with user session", func(t *testing.T) {
			event, err := mocks.Event(db, ctx)
			require.NoError(t, err)
			require.Nil(t, event.ProductID)

			countBefore, err := countRegistrations()
			require.NoError(t, err)

			url := fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID)
			resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
			assert.NoError(t, err)
			mocks.AssertRedirectResponse(t, resp.StatusCode)

			redirectedTo := resp.Header.Get("Location")
			assert.Equal(t, "/events/"+event.Slug, redirectedTo)

			countAfter, err := countRegistrations()
			require.NoError(t, err)

			assert.Equal(t, countBefore+1, countAfter)
		})

		t.Run("POST /event_registrations/:event_id is a no-op when the user is already registered", func(t *testing.T) {
			event, err := mocks.Event(db, ctx)
			require.NoError(t, err)
			require.Nil(t, event.ProductID)

			_, err = mocks.EventRegistration(db, ctx, event, user)
			require.NoError(t, err)

			countBefore, err := countRegistrations()
			require.NoError(t, err)

			url := fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID)
			resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
			assert.NoError(t, err)
			mocks.AssertRedirectResponse(t, resp.StatusCode)

			redirectedTo := resp.Header.Get("Location")
			assert.Equal(t, "/events/"+event.Slug, redirectedTo)

			countAfter, err := countRegistrations()
			require.NoError(t, err)
			assert.Equal(t, countBefore, countAfter)
		})

		t.Run("GET /events/:event_id/register signs user up for an event", func(t *testing.T) {
			event, err := mocks.Event(db, ctx)
			require.NoError(t, err)
			require.Nil(t, event.ProductID)

			countBefore, err := countRegistrations()
			require.NoError(t, err)

			resp, err := client.Get(fmt.Sprintf("%s/events/%s/register", server.URL, event.ID))
			assert.NoError(t, err)
			mocks.AssertRedirectResponse(t, resp.StatusCode)

			redirectedTo := resp.Header.Get("Location")
			assert.Equal(t, "/events/"+event.Slug, redirectedTo)

			countAfter, err := countRegistrations()
			require.NoError(t, err)
			assert.Equal(t, countBefore+1, countAfter)
		})

		t.Run("DELETE /event_registrations/:event_id unregisters user", func(t *testing.T) {
			event, err := mocks.Event(db, ctx)
			require.NoError(t, err)

			_, err = mocks.EventRegistration(db, ctx, event, user)
			require.NoError(t, err)

			countBefore, err := countRegistrations()
			require.NoError(t, err)

			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID), nil)
			resp, err := client.Do(req)
			assert.NoError(t, err)
			mocks.AssertRedirectResponse(t, resp.StatusCode)

			redirectedTo := resp.Header.Get("Location")
			assert.Equal(t, "/events/"+event.Slug, redirectedTo)

			countAfter, err := countRegistrations()
			require.NoError(t, err)
			assert.Equal(t, countBefore-1, countAfter)
		})

		t.Run("DELETE /event_registrations/:event_id is a no-op when not registered", func(t *testing.T) {
			event, err := mocks.Event(db, ctx)
			require.NoError(t, err)

			countBefore, err := countRegistrations()
			require.NoError(t, err)

			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID), nil)
			resp, err := client.Do(req)
			assert.NoError(t, err)
			mocks.AssertRedirectResponse(t, resp.StatusCode)

			countAfter, err := countRegistrations()
			require.NoError(t, err)
			assert.Equal(t, countBefore, countAfter)
		})

		t.Run("POST /event_registrations/:event_id returns 404 for paid event", func(t *testing.T) {
			event, err := mocks.PaidEvent(db, ctx)
			require.NoError(t, err)

			url := fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID)
			resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("POST /event_registrations/:event_id returns 400 for invalid event_id", func(t *testing.T) {
			url := fmt.Sprintf("%s/event_registrations/not-a-uuid", server.URL)
			resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("as anonymous user", func(t *testing.T) {
		anonClient, err := mocks.ClientWithSession(store, mustParseURL(server.URL), nil)
		require.NoError(t, err)

		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		t.Run("POST /event_registrations/:event_id redirects to sign-in", func(t *testing.T) {
			url := fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID)
			resp, err := anonClient.Post(url, "application/x-www-form-urlencoded", nil)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Location"), "/sign-in")
		})

		t.Run("GET /events/:event_id/register redirects to sign-in", func(t *testing.T) {
			resp, err := anonClient.Get(fmt.Sprintf("%s/events/%s/register", server.URL, event.ID))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Location"), "/sign-in")
		})

		t.Run("DELETE /event_registrations/:event_id redirects to sign-in", func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID), nil)
			resp, err := anonClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Location"), "/sign-in")
		})
	})
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
