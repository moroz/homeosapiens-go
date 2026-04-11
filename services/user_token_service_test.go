package services_test

import (
	"testing"
	"time"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	user, err := mocks.UniqueUser(db, t.Context())
	require.NoError(t, err)

	srv := services.NewUserTokenService(db)
	validity := 24 * 60 * 60 * time.Second

	t.Run("email verification token", func(t *testing.T) {
		token, err := srv.IssueEmailVerificationTokenForUser(t.Context(), user, validity)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.PlaintextToken)
		assert.NotEmpty(t, token.Token)
		assert.Equal(t, "email_verification", token.Context)
	})

	t.Run("access token", func(t *testing.T) {
		token, err := srv.IssueAccessTokenForUser(t.Context(), user, validity)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.Token)
		assert.Equal(t, "access", token.Context)
	})
}
