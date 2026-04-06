package workers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	appi18n "github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func workerConfig(db queries.DBTX, mailer mailers.Mailer, bundle *i18n.Bundle) *river.Workers {
	workers := river.NewWorkers()

	river.AddWorker(workers, &SendOrderConfirmationWorker{
		db:     db,
		mailer: mailer,
		bundle: bundle,
	})

	return workers
}

func Client(db queries.DBTX) (*river.Client[pgx.Tx], error) {
	conn, ok := db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("Failed to cast database connection as *pgxpool.Pool, got: %T", db)
	}

	return river.NewClient(riverpgxv5.New(conn), &river.Config{})
}

func SetupWorkers(ctx context.Context, db *pgxpool.Pool, mailer mailers.Mailer) (*river.Client[pgx.Tx], error) {
	bundle, err := appi18n.InitBundle()

	client, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workerConfig(db, mailer, bundle),
	})

	if err != nil {
		return nil, err
	}

	if err := client.Start(ctx); err != nil {
		return nil, err
	}

	return client, err
}
