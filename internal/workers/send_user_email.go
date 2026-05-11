package workers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
)

type SendUserEmailWorker struct {
	river.WorkerDefaults[jobs.SendUserEmailArgs]
	db     *pgxpool.Pool
	mailer mailers.Mailer
	bundle *i18n.Bundle
}

func (w *SendUserEmailWorker) Work(ctx context.Context, job *river.Job[jobs.SendUserEmailArgs]) error {
	user, err := queries.New(w.db).GetUserByID(ctx, job.Args.UserID)
	if err != nil {
		return err
	}

	tx, err := w.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	service := services.NewEmailVerificationService(tx)

	token, err := service.IssueEmailVerificationTokenForUser(ctx, user)
	if err != nil {
		return err
	}

	userMailer := mailers.NewUserMailer(w.mailer, w.bundle)

	if err := userMailer.SendUserEmailVerification(ctx, token); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
