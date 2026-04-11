package workers

import (
	"context"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
)

type SendUserEmailWorker struct {
	river.WorkerDefaults[jobs.SendUserEmailArgs]
	db     queries.DBTX
	mailer mailers.Mailer
	bundle *i18n.Bundle
}

func (w *SendUserEmailWorker) Work(ctx context.Context, job *river.Job[jobs.SendUserEmailArgs]) error {
	userTokenService := services.NewUserTokenService(w.db)
	user, err := queries.New(w.db).GetUserByID(ctx, job.Args.UserID)
	if err != nil {
		return err
	}

	token, err := userTokenService.IssueEmailVerificationTokenForUser(ctx, user, config.EmailVerificationTokenValidity)
	if err != nil {
		return err
	}

	userMailer := mailers.NewUserMailer(w.mailer, w.bundle)

	return userMailer.SendUserEmailVerification(ctx, token)
}
