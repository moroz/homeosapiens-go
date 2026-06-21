package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/require"
)

func buildVerifyEmailRequest(t *testing.T, email string) *http.Request {
	params := url.Values{"email": {email}}
	body := bytes.NewBufferString(params.Encode())
	req, err := http.NewRequest("POST", "/email-verifications", body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err)
	return req
}

func TestUserVerificationController_Create(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	stripeSrv := mocks.NewMockStripeService(t)

	r := router.Router(db, store, stripeSrv)

	t.Run("schedules an email job with valid params", func(t *testing.T) {
		user, err := mocks.User(db, ctx)

		require.NoError(t, err)
		require.Nil(t, user.EmailConfirmedAt)

		req := buildVerifyEmailRequest(t, user.Email.String())

		_, err = db.Exec(ctx, "truncate river_job")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		mocks.AssertRedirectResponse(t, rr.Code)

		rivertest.RequireInserted(ctx, t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})

	t.Run("does not send an email when user is verified", func(t *testing.T) {
		user, err := mocks.User(db, ctx, func(params *types.SeedUserParams) {
			params.EmailConfirmed = true
		})

		require.NoError(t, err)
		require.NotNil(t, user.EmailConfirmedAt)

		req := buildVerifyEmailRequest(t, user.Email.String())

		_, err = db.Exec(ctx, "truncate river_job")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		mocks.AssertRedirectResponse(t, rr.Code)

		rivertest.RequireNotInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})

	t.Run("does not send an email when the user does not exist", func(t *testing.T) {
		req := buildVerifyEmailRequest(t, "non-existent@example.com")

		_, err = db.Exec(ctx, "truncate river_job")
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		mocks.AssertRedirectResponse(t, rr.Code)
		rivertest.RequireNotInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})
}
