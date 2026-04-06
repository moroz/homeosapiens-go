package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/internal/workers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	mailer, err := mailers.NewSMTPMailer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	if err != nil {
		log.Fatal(err)
	}

	riverClient, err := workers.SetupWorkers(ctx, db, mailer)
	if err != nil {
		log.Fatal(err)
	}

	sessionStore, err := session.NewStore(config.SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	stripeService := services.NewStripeService(config.StripeSecretKey, config.StripeWebhookSigningSecret)

	r := router.Router(db, sessionStore, stripeService, mailer)

	// Stop the River client on SIGINT or SIGTERM
	go func() {
		<-ctx.Done()
		log.Printf("Stopping river client...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := riverClient.Stop(shutdownCtx); err != nil {
			log.Printf("River shutdown error: %v", err)
		}
	}()

	// Block the main goroutine on echo server. River workers will be stopped in the goroutine.
	sc := echo.StartConfig{
		Address:         ":" + config.AppPort,
		GracefulTimeout: 30 * time.Second,
		OnShutdownError: func(err error) {
			log.Printf("HTTP server shutdown error: %v", err)
		},
	}
	if err := sc.Start(ctx, r); err != nil {
		log.Fatal(err)
	}
}
