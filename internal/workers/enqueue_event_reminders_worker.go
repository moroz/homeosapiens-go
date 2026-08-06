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
// everybody registered for an event that is about to start. Claiming a
// registration stamps it in the same statement that selects it, so an attendee
// is only ever reminded once per lead time, even with several workers running.
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

	// The two claim queries return the same shape, but sqlc gives each its own
	// row type, so they are flattened into one list here.
	type registration struct {
		eventId, userId uuid.UUID
	}
	var registrations []registration

	switch lead {
	case jobs.ReminderLead24h:
		rows, err := q.ClaimEventRegistrationsForReminder24h(ctx)
		if err != nil {
			return err
		}
		for _, row := range rows {
			registrations = append(registrations, registration{row.EventID, row.UserID})
		}
	case jobs.ReminderLead1h:
		rows, err := q.ClaimEventRegistrationsForReminder1h(ctx)
		if err != nil {
			return err
		}
		for _, row := range rows {
			registrations = append(registrations, registration{row.EventID, row.UserID})
		}
	default:
		return fmt.Errorf("unknown reminder lead time: %q", lead)
	}

	if len(registrations) == 0 {
		return nil
	}

	client, err := jobs.NewClient(w.db)
	if err != nil {
		return err
	}

	params := make([]river.InsertManyParams, len(registrations))
	for i, r := range registrations {
		params[i] = river.InsertManyParams{
			Args: jobs.SendEventReminderEmailArgs{
				UserID:  r.userId,
				EventID: r.eventId,
				Lead:    lead,
			},
		}
	}

	if _, err := client.InsertManyTx(ctx, tx, params); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
