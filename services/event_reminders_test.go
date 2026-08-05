package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The reminder scanner relies on the claim queries stamping the events they
// return, so that an event is reminded about exactly once per lead time no
// matter how often the scan runs.
func TestClaimEventsForReminder(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, products cascade")
	require.NoError(t, err)

	now := time.Now().UTC()

	// Timestamps are compared against now() in the database, so they are built in
	// UTC rather than in the local zone.
	inSixHours, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
		p.StartsAt = now.Add(6 * time.Hour)
		p.EndsAt = now.Add(8 * time.Hour)
	})
	require.NoError(t, err)

	inHalfHour, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
		p.StartsAt = now.Add(30 * time.Minute)
		p.EndsAt = now.Add(90 * time.Minute)
	})
	require.NoError(t, err)

	// Too far out for either reminder.
	_, err = mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
		p.StartsAt = now.Add(72 * time.Hour)
		p.EndsAt = now.Add(74 * time.Hour)
	})
	require.NoError(t, err)

	// Unpublished events are never reminded about.
	_, err = mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
		p.StartsAt = now.Add(6 * time.Hour)
		p.EndsAt = now.Add(8 * time.Hour)
		p.Published = false
	})
	require.NoError(t, err)

	q := queries.New(db)

	claimed24h, err := q.ClaimEventsForReminder24h(ctx)
	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{inSixHours.ID}, claimed24h)

	claimed1h, err := q.ClaimEventsForReminder1h(ctx)
	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{inHalfHour.ID}, claimed1h)

	// A second scan has nothing left to claim.
	claimed24h, err = q.ClaimEventsForReminder24h(ctx)
	require.NoError(t, err)
	assert.Empty(t, claimed24h)

	claimed1h, err = q.ClaimEventsForReminder1h(ctx)
	require.NoError(t, err)
	assert.Empty(t, claimed1h)
}
