package workers

import (
	"context"
	"fmt"

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
	userMailer := mailers.NewUserMailer(w.mailer, w.bundle)

	switch job.Args.EmailType {
	case jobs.UserEmailTypeEmailVerification:
		service := services.NewEmailVerificationService(w.db)

		token, err := service.IssueEmailVerificationTokenForUser(ctx, user)
		if err != nil {
			return err
		}
		return userMailer.SendUserEmailVerification(ctx, token)

	case jobs.UserEmailTypePasswordReset:
		service := services.NewUserPasswordResetService(w.db)
		token, err := service.IssuePasswordResetTokenForUser(ctx, user)
		if err != nil {
			return err
		}
		return userMailer.SendUserPasswordResetEmail(ctx, token)
	default:
		return fmt.Errorf("Unknown email type %v", job.Args.EmailType)
	}
}
