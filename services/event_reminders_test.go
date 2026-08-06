package services_test

import (
	"testing"
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The reminder scanner relies on the claim queries stamping the registrations
// they return, so that an attendee is reminded exactly once per lead time no
// matter how often the scan runs.
func TestClaimEventRegistrationsForReminder(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, products, event_registrations cascade")
	require.NoError(t, err)

	now := time.Now().UTC()

	// Timestamps are compared against now() in the database, so they are built in
	// UTC rather than in the local zone.
	newEvent := func(startsIn time.Duration, overrides ...func(*queries.UpsertEventParams)) *queries.Event {
		event, err := mocks.Event(db, ctx, append([]func(*queries.UpsertEventParams){
			func(p *queries.UpsertEventParams) {
				p.StartsAt = now.Add(startsIn)
				p.EndsAt = now.Add(startsIn + 2*time.Hour)
			},
		}, overrides...)...)
		require.NoError(t, err)
		return event
	}

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	inSixHours := newEvent(6 * time.Hour)
	inHalfHour := newEvent(30 * time.Minute)
	// Too far out for either reminder, and an unpublished event is never
	// reminded about at all.
	farOut := newEvent(72 * time.Hour)
	unpublished := newEvent(6*time.Hour, func(p *queries.UpsertEventParams) { p.Published = false })

	for _, event := range []*queries.Event{inSixHours, inHalfHour, farOut, unpublished} {
		_, err := mocks.EventRegistration(db, ctx, event, user)
		require.NoError(t, err)
	}

	q := queries.New(db)

	claimed24h, err := q.ClaimEventRegistrationsForReminder24h(ctx)
	require.NoError(t, err)
	require.Len(t, claimed24h, 1)
	assert.Equal(t, inSixHours.ID, claimed24h[0].EventID)
	assert.Equal(t, user.ID, claimed24h[0].UserID)

	claimed1h, err := q.ClaimEventRegistrationsForReminder1h(ctx)
	require.NoError(t, err)
	require.Len(t, claimed1h, 1)
	assert.Equal(t, inHalfHour.ID, claimed1h[0].EventID)

	// A second scan has nothing left to claim.
	claimed24h, err = q.ClaimEventRegistrationsForReminder24h(ctx)
	require.NoError(t, err)
	assert.Empty(t, claimed24h)

	claimed1h, err = q.ClaimEventRegistrationsForReminder1h(ctx)
	require.NoError(t, err)
	assert.Empty(t, claimed1h)

	// Somebody registering after that scan still gets the reminder: the stamp is
	// theirs, not the event's.
	latecomer, err := mocks.User(db, ctx)
	require.NoError(t, err)
	_, err = mocks.EventRegistration(db, ctx, inSixHours, latecomer)
	require.NoError(t, err)

	claimed24h, err = q.ClaimEventRegistrationsForReminder24h(ctx)
	require.NoError(t, err)
	require.Len(t, claimed24h, 1)
	assert.Equal(t, latecomer.ID, claimed24h[0].UserID)
}
