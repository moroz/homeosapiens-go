package services_test

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

	t.Run("returns validation error on duplicate email address", func(t *testing.T) {
		user, err := srv.RegisterUser(t.Context(), params)
		assert.Nil(t, user)

		validationErrors, ok := errors.AsType[validation.Errors](err)
		assert.True(t, ok)

		assert.NotEmpty(t, validationErrors["email"])
	})
}
