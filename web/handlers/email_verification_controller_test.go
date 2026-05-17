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

func TestUserVerificationController_Create(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	stripeSrv := mocks.NewMockStripeService(t)

	user, err := services.NewUserService(db).CreateUser(ctx, &types.SeedUserParams{
		GivenName:  "Test",
		FamilyName: "User",
		Email:      mocks.UniqueEmail(),
		Country:    "US",
		Password:   "foobar",
	})

	require.NoError(t, err)
	require.Nil(t, user.EmailConfirmedAt)

	r := router.Router(db, store, stripeSrv)

	params := url.Values{"email": {user.Email.String()}}
	body := bytes.NewBufferString(params.Encode())
	req, err := http.NewRequest("POST", "/email-verifications", body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.GreaterOrEqual(t, rr.Code, 300)
	assert.Less(t, rr.Code, 400)
}
