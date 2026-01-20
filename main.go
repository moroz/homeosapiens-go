package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/handlers"
	"github.com/moroz/homeosapiens-go/i18n"
	"github.com/moroz/securecookie"
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

	sessionStore, err := securecookie.NewStore(config.SessionKey)

	r := handlers.EchoRouter(db, bundle, sessionStore)
	listenOn := ":" + config.AppPort
	log.Printf("Listening on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, r))
}
