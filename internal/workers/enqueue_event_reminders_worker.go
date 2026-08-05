package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/riverqueue/river"
)

// EnqueueEventRemindersWorker runs on a schedule and fans reminders out to
// everybody registered for an event that is about to start. Claiming an event
// stamps it in the same statement that selects it, so an event is only ever
// reminded about once per lead time, even with several workers running.
type EnqueueEventRemindersWorker struct {
	river.WorkerDefaults[jobs.EnqueueEventRemindersArgs]
	db *pgxpool.Pool
}

func (w *EnqueueEventRemindersWorker) Work(ctx context.Context, job *river.Job[jobs.EnqueueEventRemindersArgs]) error {
	for _, lead := range []jobs.ReminderLead{jobs.ReminderLead24h, jobs.ReminderLead1h} {
		if err := w.enqueueForLead(ctx, lead); err != nil {
			return fmt.Errorf("EnqueueEventReminders (%s): %w", lead, err)
		}
	}

	return nil
}

func (w *EnqueueEventRemindersWorker) enqueueForLead(ctx context.Context, lead jobs.ReminderLead) error {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := queries.New(tx)

	var eventIds []uuid.UUID
	switch lead {
	case jobs.ReminderLead24h:
		eventIds, err = q.ClaimEventsForReminder24h(ctx)
	case jobs.ReminderLead1h:
		eventIds, err = q.ClaimEventsForReminder1h(ctx)
	default:
		return fmt.Errorf("unknown reminder lead time: %q", lead)
	}
	if err != nil {
		return err
	}

	if len(eventIds) == 0 {
		return nil
	}

	registrations, err := q.ListUserIDsForEventRegistrations(ctx, eventIds)
	if err != nil {
		return err
	}

	if len(registrations) == 0 {
		// The events stay stamped: nobody was registered, so there is nothing to
		// remind them about later either.
		return tx.Commit(ctx)
	}

	client, err := jobs.NewClient(w.db)
	if err != nil {
		return err
	}

	params := make([]river.InsertManyParams, len(registrations))
	for i, registration := range registrations {
		params[i] = river.InsertManyParams{
			Args: jobs.SendEventReminderEmailArgs{
				UserID:  registration.UserID,
				EventID: registration.EventID,
				Lead:    lead,
			},
		}
	}

	if _, err := client.InsertManyTx(ctx, tx, params); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
