package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetFlow(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	stripeSrv := mocks.NewMockStripeService(t)
	r := router.Router(db, store, stripeSrv)

	t.Run("renders a form", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/reset-password", nil)
		require.NoError(t, err)

		tt := httptest.NewRecorder()

		r.ServeHTTP(tt, req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, tt.Code)
	})
}
