package services_test

import (
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPasswordResetService(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	srv := services.NewUserPasswordResetService(db)
	user, err := mocks.UniqueUser(db, ctx)
	require.NoError(t, err)

	t.Run("MaybeIssuePasswordResetTokenForUser", func(t *testing.T) {
		_, err = db.Exec(ctx, "truncate river_job")
		require.NoError(t, err)

		sent, err := srv.MaybeIssuePasswordResetTokenForUser(ctx, user.Email.String())
		assert.NoError(t, err)
		assert.True(t, sent)

		rivertest.RequireInserted(ctx, t, riverpgxv5.New(db), jobs.SendUserEmailArgs{}, nil)
	})

	t.Run("UpdateUserPassword", func(t *testing.T) {
		token, err := srv.IssuePasswordResetTokenForUser(ctx, user)
		require.NoError(t, err)
		require.NotNil(t, token)

		updated, err := srv.UpdateUserPassword(ctx, token.PlaintextToken, &types.UpdateUserPasswordRequest{
			Password:             "new_password",
			PasswordConfirmation: "new_password",
		})
		assert.NoError(t, err)
		assert.NotEqual(t, user.PasswordHash, updated.PasswordHash)

		valid, err := argon2id.ComparePasswordAndHash("new_password", *updated.PasswordHash)
		assert.True(t, valid)
		assert.NoError(t, err)
	})
}
