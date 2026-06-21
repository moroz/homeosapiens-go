package handlers_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRegistrationController_Create(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)

	_, err = db.Exec(ctx, "truncate event_registrations")
	require.NoError(t, err)

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	r := router.Router(db, store, mocks.NewMockStripeService(t))

	server := httptest.NewServer(r)
	defer server.Close()

	client, err := mocks.ClientWithUser(&mocks.ClientWithUserInput{
		Store:   store,
		Server:  server,
		User:    user,
		DB:      db,
		Context: ctx,
	})
	require.NoError(t, err)

	countRegistrations := func() (int, error) {
		var count int
		err := db.QueryRow(ctx, "select count(*) from event_registrations").Scan(&count)
		return count, err
	}

	t.Run("POST /event_registrations/:event_id signs up for free event with user session", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)
		require.Nil(t, event.ProductID)

		countBefore, err := countRegistrations()
		require.NoError(t, err)

		url := fmt.Sprintf("%s/event_registrations/%s", server.URL, event.ID)
		resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, resp.StatusCode, 300)
		assert.Less(t, resp.StatusCode, 400)

		countAfter, err := countRegistrations()
		require.NoError(t, err)

		assert.Equal(t, countBefore+1, countAfter)
	})
}
