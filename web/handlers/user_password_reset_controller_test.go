package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildResetPasswordRequest(t *testing.T, email string) *http.Request {
	params := url.Values{"email": {email}}
	body := bytes.NewBufferString(params.Encode())
	req, err := http.NewRequest("POST", "/reset-password", body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err)
	return req
}

func TestPasswordResetFlow(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	stripeSrv := mocks.NewMockStripeService(t)
	r := router.Router(db, store, stripeSrv)

	t.Run("GET /forgot-password renders a form", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/reset-password", nil)
		require.NoError(t, err)

		tt := httptest.NewRecorder()

		r.ServeHTTP(tt, req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, tt.Code)
	})

	t.Run("POST /forgot-password with valid params", func(t *testing.T) {
		user, err := services.NewUserService(db).CreateUser(ctx, &types.SeedUserParams{
			Email:    mocks.UniqueEmail(),
			Password: "foobar",
		})
		require.NoError(t, err)

		req := buildResetPasswordRequest(t, user.Email.String())

		tt := httptest.NewRecorder()

		r.ServeHTTP(tt, req)
		assert.NoError(t, err)
	})
}
