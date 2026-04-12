package services_test

import (
	"testing"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)

	defer db.Close()

	db.Exec(t.Context(), "truncate users cascade")

	params := &types.RegisterUserParams{
		PreferredLocale:      "en",
		GivenName:            "John",
		FamilyName:           "Smith",
		Email:                mocks.UniqueEmail(),
		Password:             "foobar2000",
		PasswordConfirmation: "foobar2000",
	}

	err = params.Validate()
	require.NoError(t, err)

	srv := services.NewUserRegistrationService(db)

	t.Run("registers user with valid params", func(t *testing.T) {
		db.Exec(t.Context(), "truncate river_job")

		user, err := srv.RegisterUser(t.Context(), params)
		assert.NoError(t, err)
		assert.NotNil(t, user)

		rivertest.RequireInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendUserEmailArgs{}, nil)
	})
}

func TestVerifyEmailAddress(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	srv := services.NewUserRegistrationService(db)

	t.Run("marks email verified with valid token", func(t *testing.T) {
		user, err := mocks.UniqueUser(db, t.Context(), func(params *queries.UpsertUserFromSeedDataParams) {
			params.EmailConfirmedAt = nil
		})
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
		user, err := mocks.UniqueUser(db, t.Context(), func(params *queries.UpsertUserFromSeedDataParams) {
			params.EmailConfirmedAt = new(time.Now())
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
		user, err := mocks.UniqueUser(db, t.Context(), func(params *queries.UpsertUserFromSeedDataParams) {
			params.EmailConfirmedAt = nil
		})
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
