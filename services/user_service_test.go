package services_test

import (
	"context"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/bincyber/go-sqlcrypter"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (queries.DBTX, error) {
	connString := config.MustGetenv("TEST_DATABASE_URL")
	return pgxpool.New(ctx, connString)
}

func TestAuthenticateUserByEmailPassword(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)

	db.Exec(t.Context(), "truncate users cascade")

	hash, err := argon2id.CreateHash("foobar", argon2id.DefaultParams)
	require.NoError(t, err)

	user, err := queries.New(db).InsertUser(t.Context(), &queries.InsertUserParams{
		Email:        sqlcrypter.NewEncryptedBytes("user@example.com"),
		EmailHash:    crypto.HashEmail("user@example.com"),
		GivenName:    sqlcrypter.NewEncryptedBytes("Example"),
		FamilyName:   sqlcrypter.NewEncryptedBytes("User"),
		PasswordHash: &hash,
	})
	require.NoError(t, err)
	assert.NotNil(t, user)

	subject := services.NewUserService(db)

	t.Run("with valid params", func(t *testing.T) {
		emails := []string{
			"user@example.com",
			"USER@example.com",
			"USER@EXAMPLE.COM",
		}

		for _, email := range emails {
			actual, err := subject.AuthenticateUserByEmailPassword(t.Context(), email, "foobar")
			assert.NoError(t, err)
			assert.Equal(t, user.ID, actual.ID)
		}
	})

	t.Run("with invalid password", func(t *testing.T) {
		passwords := []string{
			"",
			"Foobar",
			"foobar ",
			"foo bar",
			"FOOBAR",
		}

		for _, pass := range passwords {
			actual, err := subject.AuthenticateUserByEmailPassword(t.Context(), user.Email.String(), pass)
			assert.Nil(t, actual)
			assert.ErrorIs(t, err, services.ErrInvalidPassword)
		}
	})
}
