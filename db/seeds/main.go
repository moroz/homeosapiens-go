package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	users := []*types.SeedUserParams{
		{
			Email:          "karol@moroz.dev",
			GivenName:      "Karol",
			FamilyName:     "Moroz",
			Country:        "PL",
			Role:           queries.UserRoleAdministrator,
			EmailConfirmed: true,
		},
	}

	userService := services.NewUserService(db)
	for _, user := range users {
		if _, err := userService.CreateUser(context.Background(), user); err != nil {
			log.Fatal(err)
		}
	}
}
