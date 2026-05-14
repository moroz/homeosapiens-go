package services_test

import (
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyEmailAddress(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	srv := services.NewEmailVerificationService(db)

	t.Run("marks email verified with valid token", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context())
		require.NoError(t, err)
		require.Nil(t, user.EmailConfirmedAt)

		token, err := services.NewUserTokenService(db).IssueEmailVerificationTokenForUser(t.Context(), user, config.EmailVerificationTokenValidity)
		require.NoError(t, err)

		var exists bool
		err = db.QueryRow(t.Context(), "select exists (select from user_tokens where token = $1)", crypto.HashUserToken(token.PlaintextToken)).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)

		user, err = srv.VerifyEmailAddress(t.Context(), token.EncodeToken())
		assert.NoError(t, err)
		assert.NotNil(t, user.EmailConfirmedAt)

		err = db.QueryRow(t.Context(), "select exists (select from user_tokens where token = $1)", crypto.HashUserToken(token.PlaintextToken)).Scan(&exists)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("does nothing when a user is already confirmed", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context(), func(params *types.SeedUserParams) {
			params.EmailConfirmed = true
		})
		require.NoError(t, err)
		require.NotNil(t, user.EmailConfirmedAt)

		token, err := services.NewUserTokenService(db).IssueEmailVerificationTokenForUser(t.Context(), user, config.EmailVerificationTokenValidity)
		require.NoError(t, err)

		var exists bool
		err = db.QueryRow(t.Context(), "select exists (select from user_tokens where token = $1)", crypto.HashUserToken(token.PlaintextToken)).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)

		actual, err := srv.VerifyEmailAddress(t.Context(), token.EncodeToken())
		assert.ErrorIs(t, err, services.ErrTokenInvalid)
		assert.Nil(t, actual)

		updatedUser, err := queries.New(db).GetUserByID(t.Context(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.EmailConfirmedAt, updatedUser.EmailConfirmedAt)
		assert.Equal(t, user.UpdatedAt, updatedUser.UpdatedAt)
	})

	t.Run("returns error and does nothing when a token is expired", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context())
		require.NoError(t, err)
		require.Nil(t, user.EmailConfirmedAt)

		token, err := services.NewUserTokenService(db).IssueEmailVerificationTokenForUser(t.Context(), user, config.EmailVerificationTokenValidity)
		require.NoError(t, err)

		_, err = db.Exec(t.Context(), "update user_tokens set inserted_at = now() - interval '5 days', valid_until = now() - interval '1 day' where id = $1", token.UserToken.ID)
		require.NoError(t, err)

		actual, err := srv.VerifyEmailAddress(t.Context(), token.EncodeToken())
		assert.ErrorIs(t, err, services.ErrTokenInvalid)
		assert.Nil(t, actual)

		user, err = queries.New(db).GetUserByID(t.Context(), user.ID)
		require.NoError(t, err)
		assert.Nil(t, user.EmailConfirmedAt)
	})
}
