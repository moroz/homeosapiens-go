package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
)

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	sessionStore, err := session.NewStore(config.SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	stripeService := services.NewStripeService(config.StripeSecretKey, config.StripeWebhookSigningSecret)

	mailer, err := mailers.NewSMTPMailer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	if err != nil {
		log.Fatal(err)
	}

	r := router.Router(db, sessionStore, stripeService, mailer)
	listenOn := ":" + config.AppPort
	log.Fatal(r.Start(listenOn))
}
