package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

func main() {
	fmt.Println(config.DatabaseUrl)
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(context.Background())

	log.Printf("Cleaning database...")
	_, err = db.Exec(context.Background(), "truncate events, hosts, assets, venues, events_hosts, event_prices, event_registrations, user_tokens, videos, video_sources")
	if err != nil {
		log.Fatal(err)
	}

	users := []*types.CreateUserParams{
		{
			Email:      "karol@moroz.dev",
			GivenName:  "Karol",
			FamilyName: "Moroz",
			Country:    "PL",
			Password:   "foobar",
			Role:       queries.UserRoleAdministrator,
		},
		{
			GivenName:  "Sanjay",
			FamilyName: "Modi",
			Email:      "sanjay.modi@example.com",
			Country:    "IN",
			Password:   "foobar",
			Role:       queries.UserRoleRegular,
		},
	}

	log.Printf("Creating users...")
	userService := services.NewUserService(db)
	for _, user := range users {
		if _, err := userService.CreateUser(context.Background(), user); err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit(context.Background())
}
