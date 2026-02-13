package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost/postgres?sslmode=disable&options=-c%20TimeZone%3DAmerica%2FNew_York"

	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	var t string
	err = conn.QueryRow(context.Background(), "select now()::text").Scan(&t)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(t)
}
