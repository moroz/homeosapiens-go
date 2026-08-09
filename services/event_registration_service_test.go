package services_test

import (
	"errors"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRegistrationService_AdminCreateEventRegistration(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, event_registrations, river_job cascade")
	require.NoError(t, err)

	srv := services.NewEventRegistrationService(db)

	t.Run("registers a user for a published, ongoing event", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		registration, err := srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: event.ID,
			UserID:  user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, registration)
		assert.Equal(t, event.ID, registration.EventID)
		assert.Equal(t, user.ID, registration.UserID)
	})

	t.Run("rejects an unpublished event", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.Published = false
		})
		require.NoError(t, err)
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		_, err = srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: event.ID,
			UserID:  user.ID,
		})

		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Contains(t, verrs, "eventId")
	})

	t.Run("rejects an event that has already ended", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.StartsAt = time.Now().UTC().Add(-2 * time.Hour)
			p.EndsAt = time.Now().UTC().Add(-time.Hour)
		})
		require.NoError(t, err)
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		_, err = srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: event.ID,
			UserID:  user.ID,
		})

		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Contains(t, verrs, "eventId")
	})

	t.Run("rejects an unknown event", func(t *testing.T) {
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		_, err = srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: uuid.Must(uuid.NewV7()),
			UserID:  user.ID,
		})
		assert.Error(t, err)
	})

	t.Run("rejects an unknown user", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		_, err = srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: event.ID,
			UserID:  uuid.Must(uuid.NewV7()),
		})
		assert.Error(t, err)
	})

	t.Run("rejects a user who is already registered", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)
		user, err := mocks.User(db, ctx)
		require.NoError(t, err)

		_, err = mocks.EventRegistration(db, ctx, event, user)
		require.NoError(t, err)

		_, err = srv.AdminCreateEventRegistration(ctx, &types.EnrollStudentForEventInput{
			EventID: event.ID,
			UserID:  user.ID,
		})

		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Contains(t, verrs, "userId")
	})
}
