package workers

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
)

type SendOrderConfirmationArgs struct {
	OrderID uuid.UUID `json:"order_id"`
}

func (SendOrderConfirmationArgs) Kind() string { return "SendOrderConfirmation" }

type SendOrderConfirmationWorker struct {
	river.WorkerDefaults[SendOrderConfirmationArgs]
	db     queries.DBTX
	mailer mailers.Mailer
	bundle *i18n.Bundle
}

func (w *SendOrderConfirmationWorker) Work(ctx context.Context, job *river.Job[SendOrderConfirmationArgs]) error {
	srv := services.NewOrderRepository(w.db)
	order, err := srv.GetOrderDetails(ctx, job.Args.OrderID)
	if err != nil {
		return err
	}

	orderMailer := mailers.NewOrderMailer(w.mailer, w.bundle)

	return orderMailer.SendOrderConfirmation(ctx, order)
}
