package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/handlers"
)

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.Router(db)
	log.Fatal(http.ListenAndServe(":3000", r))
}
