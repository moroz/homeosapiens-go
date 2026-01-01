package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

func MustParseUUID(u string) pgtype.UUID {
	parsed := uuid.MustParse(u)
	return pgtype.UUID{
		Valid: true,
		Bytes: parsed,
	}
}

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

	assets := []*types.CreateAssetParams{
		{
			ID:               MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8b"),
			ObjectKey:        "cm7uqj3q500mglz8z2dqy8sdz.webp",
			OriginalFilename: "cm7uqj3q500mglz8z2dqy8sdz.webp",
		},
		{
			ID:               MustParseUUID("019b0c7c-c3c4-71c3-a630-7b33a847ca2a"),
			ObjectKey:        "019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg",
			OriginalFilename: "019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg",
		},
	}
	log.Printf("Creating assets...")
	for _, asset := range assets {
		params := &queries.UpsertAssetParams{
			ID:               asset.ID,
			ObjectKey:        asset.ObjectKey,
			OriginalFilename: &asset.OriginalFilename,
		}
		if _, err := queries.New(db).UpsertAsset(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	hosts := []*types.CreateHostParams{
		{
			ID:               MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8a"),
			Salutation:       "common.hosts.salutation.dr",
			GivenName:        "Sanjay",
			FamilyName:       "Modi",
			ProfilePictureId: MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8b"),
			Country:          "IN",
		},
		{
			ID:               MustParseUUID("019b0c71-fde2-76b7-8c71-21c2e9ea23a5"),
			Salutation:       "common.hosts.salutation.dr",
			GivenName:        "Herman",
			FamilyName:       "Jeggels",
			ProfilePictureId: MustParseUUID("019b0c7c-c3c4-71c3-a630-7b33a847ca2a"),
			Country:          "ZA",
		},
	}
	log.Printf("Creating hosts...")
	for _, host := range hosts {
		params := &queries.UpsertHostParams{
			ID:               host.ID,
			Salutation:       &host.Salutation,
			GivenName:        host.GivenName,
			FamilyName:       host.FamilyName,
			ProfilePictureID: host.ProfilePictureId,
			Country:          &host.Country,
		}
		if _, err := queries.New(db).UpsertHost(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit(context.Background())
}
