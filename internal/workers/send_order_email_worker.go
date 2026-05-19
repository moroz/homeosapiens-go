package workers

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
)

type SendOrderEmailWorker struct {
	river.WorkerDefaults[jobs.SendOrderEmailArgs]
	db     queries.DBTX
	mailer mailers.Mailer
	bundle *i18n.Bundle
}

func (w *SendOrderEmailWorker) Work(ctx context.Context, job *river.Job[jobs.SendOrderEmailArgs]) error {
	srv := services.NewOrderRepository(w.db)
	order, err := srv.GetOrderDetails(ctx, job.Args.OrderID)
	if err != nil {
		return err
	}

	orderMailer := mailers.NewOrderMailer(w.mailer, w.bundle)

	switch job.Args.EmailType {
	case jobs.OrderEmailTypeOrderConfirmation:
		return orderMailer.SendOrderConfirmation(ctx, order)
	case jobs.OrderEmailTypePaymentConfirmation:
		return orderMailer.SendPaymentConfirmation(ctx, order)
	default:
		return fmt.Errorf("invalid EmailType: %v", job.Args.EmailType)
	}
}
