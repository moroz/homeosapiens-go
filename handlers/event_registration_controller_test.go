package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/handlers"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (queries.DBTX, error) {
	connString := config.MustGetenv("TEST_DATABASE_URL")
	return pgxpool.New(ctx, connString)
}

func initRouter(db queries.DBTX) (http.Handler, error) {
	bundle, err := i18n.InitBundle()
	if err != nil {
		return nil, err
	}

	sessionStore, err := securecookie.NewStore(config.SessionKey)

	return handlers.Router(db, bundle, sessionStore), nil
}

var encoder = schema.NewEncoder()

func TestRegisterForEvent(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)

	bundle, err := i18n.InitBundle()
	require.NoError(t, err)

	sessionStore, err := securecookie.NewStore(config.SessionKey)
	require.NoError(t, err)

	router := handlers.Router(db, bundle, sessionStore)

	db.Exec(t.Context(), "truncate users cascade")

	eventID := "019b0c80-a410-7728-ab6b-c1eff529dfd1"
	params := types.CreateEventRegistrationParams{
		EventID:    eventID,
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john.doe@gmail.com",
		Country:    "US",
	}
	form := url.Values{}
	err = encoder.Encode(params, form)
	require.NoError(t, err)

	t.Run("when not signed in", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/event_registrations", bytes.NewBufferString(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		router.ServeHTTP(w, req)

		assert.GreaterOrEqual(t, w.Code, 200)
		assert.Less(t, w.Code, 400)
	})
}
