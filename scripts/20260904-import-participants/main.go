package main

import (
	"context"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/crypto"
)

//go:embed event_id_map.json
var EventMappingJSON []byte

// LegacyEvent maps a legacy EzyCourse event to the records that replace it:
// the event itself, the product sold for it, and the video group (playlist)
// with its recordings, if any of those exist.
type LegacyEvent struct {
	Title               string     `json:"title"`
	EventID             uuid.UUID  `json:"event_id"`
	EventProductID      *uuid.UUID `json:"event_product_id"`
	VideoGroupID        *uuid.UUID `json:"video_group_id"`
	VideoGroupProductID *uuid.UUID `json:"video_group_product_id"`
	Note                string     `json:"note,omitempty"`
}

const insertParticipantQuery = `
insert into event_registrations as er (event_id, user_id, inserted_at)
select $1, u.id, coalesce($3, now())
from users u
where u.email_hash = $2
on conflict (event_id, user_id) do update set user_id = excluded.user_id
returning er.user_id, er.inserted_at;
`

const grantProductAccessQuery = `
insert into user_product_access (user_id, product_id, granted_by_user_id, inserted_at)
values ($1, $2, $3, $4)
on conflict (user_id, product_id) do nothing;
`

const getAdminIdQuery = `
select id from users where email_hash = $1 and user_role = 'Administrator';
`

func main() {
	grantedBy := flag.String(
		"granted-by",
		"",
		"Email address of the administrator to record as the grantor of product access.",
	)
	flag.Parse()

	var mapping map[string]LegacyEvent
	if err := json.Unmarshal(EventMappingJSON, &mapping); err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(os.Stdin)

	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())

	var grantedByUserId *uuid.UUID
	if *grantedBy != "" {
		var id uuid.UUID
		if err := tx.QueryRow(
			context.Background(),
			getAdminIdQuery,
			crypto.HashEmail(*grantedBy),
		).Scan(&id); err != nil {
			log.Fatalf("Failed to find an administrator with email %v: %s", *grantedBy, err)
		}
		grantedByUserId = &id
	}

	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		event, ok := mapping[row[0]]
		if !ok {
			log.Fatalf("Failed to match event %v to new event ID!", row[0])
		}

		email := row[1]

		var insertedAt *time.Time
		if row[2] != "" {
			parsed, err := time.Parse(time.RFC3339, row[2])
			if err == nil {
				insertedAt = &parsed
			}
		}

		// The registration timestamp comes back from the insert, so that access
		// granted below carries the same timestamp as the registration itself,
		// including when the CSV had none and when the registration already
		// existed with a timestamp of its own.
		var userId uuid.UUID
		var registeredAt time.Time
		err = tx.QueryRow(
			context.Background(),
			insertParticipantQuery,
			event.EventID,
			crypto.HashEmail(email),
			insertedAt,
		).Scan(&userId, &registeredAt)

		if err != nil {
			log.Fatalf("Failed to insert user %v: %s", email, err)
		}

		// Registrants of paid events keep access to the event product and to
		// the recordings sold as a separate video group product.
		for _, productId := range []*uuid.UUID{event.EventProductID, event.VideoGroupProductID} {
			if productId == nil {
				continue
			}

			if _, err := tx.Exec(
				context.Background(),
				grantProductAccessQuery,
				userId,
				productId,
				grantedByUserId,
				registeredAt,
			); err != nil {
				log.Fatalf("Failed to grant access to product %v for user %v: %s", productId, email, err)
			}
		}
	}

	tx.Commit(context.Background())
}
