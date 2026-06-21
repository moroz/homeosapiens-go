package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexedwards/argon2id"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
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

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	stripeSrv := mocks.NewMockStripeService(t)
	r := router.Router(db, store, stripeSrv)

	t.Run("GET /forgot-password renders a form", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/reset-password", nil)
		require.NoError(t, err)
		tt := httptest.NewRecorder()
		r.ServeHTTP(tt, req)
		assert.Equal(t, http.StatusOK, tt.Code)
	})

	t.Run("POST /forgot-password with valid params", func(t *testing.T) {
		user, err := mocks.User(db, t.Context())
		require.NoError(t, err)
		db.Exec(t.Context(), "truncate river_job")

		req := buildResetPasswordRequest(t, user.Email.String())
		tt := httptest.NewRecorder()
		r.ServeHTTP(tt, req)

		assert.Equal(t, http.StatusFound, tt.Code)
		assert.Equal(t, "/sign-in", tt.Header().Get("Location"))

		rivertest.RequireInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})

	t.Run("POST /forgot-password with non-existent email", func(t *testing.T) {
		db.Exec(t.Context(), "truncate river_job")
		req := buildResetPasswordRequest(t, mocks.UniqueEmail())
		tt := httptest.NewRecorder()
		r.ServeHTTP(tt, req)
		assert.Equal(t, http.StatusFound, tt.Code)
		assert.Equal(t, "/sign-in", tt.Header().Get("Location"))
		rivertest.RequireNotInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})

	t.Run("GET /forgot-password/:token with valid params", func(t *testing.T) {
		user, err := mocks.User(db, t.Context())
		require.NoError(t, err)

		token, err := services.NewUserPasswordResetService(db).IssuePasswordResetTokenForUser(t.Context(), user)
		require.NoError(t, err)
		require.NotNil(t, token)

		req, err := http.NewRequest("GET", token.ResetPasswordPath(), nil)
		require.NoError(t, err)

		tt := httptest.NewRecorder()
		r.ServeHTTP(tt, req)

		assert.Equal(t, http.StatusOK, tt.Code)

		body, err := goquery.NewDocumentFromReader(tt.Body)
		assert.NoError(t, err)

		assert.NotEmpty(t, body.Find("form[action]").Nodes)
		assert.NotEmpty(t, body.Find("form input[name=password]").Nodes)
		assert.NotEmpty(t, body.Find("form input[name=password_confirmation]").Nodes)
	})

	t.Run("PUT /forgot-password/:token with valid params", func(t *testing.T) {
		user, err := mocks.User(db, t.Context())
		require.NoError(t, err)

		token, err := services.NewUserPasswordResetService(db).IssuePasswordResetTokenForUser(t.Context(), user)
		require.NoError(t, err)
		require.NotNil(t, token)

		password := "changed_password"
		body := url.Values{
			"password":              {password},
			"password_confirmation": {password},
		}.Encode()

		req, err := http.NewRequest("PUT", token.ResetPasswordPath(), bytes.NewReader([]byte(body)))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		tt := httptest.NewRecorder()
		r.ServeHTTP(tt, req)

		assert.GreaterOrEqual(t, tt.Code, 300)
		assert.Less(t, tt.Code, 400)

		updated, err := queries.New(db).GetUserByID(t.Context(), user.ID)
		require.NoError(t, err)
		assert.NotEqual(t, user.PasswordHash, updated.PasswordHash)

		match, _, err := argon2id.CheckHash(password, *updated.PasswordHash)
		assert.NoError(t, err)
		assert.True(t, match)
	})
}
