package workers

import (
	"context"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
)

type SendEventReminderEmailWorker struct {
	river.WorkerDefaults[jobs.SendEventReminderEmailArgs]
	db     queries.DBTX
	mailer mailers.Mailer
	bundle *i18n.Bundle
}

func (w *SendEventReminderEmailWorker) Work(ctx context.Context, job *river.Job[jobs.SendEventReminderEmailArgs]) error {
	q := queries.New(w.db)

	user, err := q.GetUserByID(ctx, job.Args.UserID)
	if err != nil {
		return err
	}

	event, err := q.GetEventById(ctx, job.Args.EventID)
	if err != nil {
		return err
	}

	data := &types.EventReminderEmailDTO{
		EventRegistrationEmailDTO: &types.EventRegistrationEmailDTO{
			Event: event,
			User:  user,
		},
		LeadKey: string(job.Args.Lead),
	}

	return mailers.NewEventMailer(w.mailer, w.bundle).SendEventReminder(ctx, data)
}
