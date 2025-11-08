package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/handlers"
	"github.com/moroz/homeosapiens-go/i18n"
)

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	bundle, err := i18n.InitBundle()
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.Router(db, bundle)
	log.Fatal(http.ListenAndServe(":"+config.AppPort, r))
}
