package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/securecookie"
)

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	sessionStore, err := securecookie.NewStore(config.SessionKey)

	r := router.Router(db, sessionStore)
	listenOn := ":" + config.AppPort
	log.Fatal(r.Start(listenOn))
}
