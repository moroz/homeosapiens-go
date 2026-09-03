package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	_ "embed"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/crypto"
)

//go:embed country-mapping.json
var CountryMapping []byte

func capitalize(name string) string {
	segs := strings.Fields(name)
	if len(segs) > 1 {
		for i, seg := range segs {
			segs[i] = capitalize(seg)
		}
		return strings.Join(segs, " ")
	}

	segs = strings.Split(name, "-")
	if len(segs) > 1 {
		for i, seg := range segs {
			segs[i] = capitalize(seg)
		}
		return strings.Join(segs, "-")
	}

	runes := []rune(strings.TrimSpace(name))

	for i, r := range runes {
		if i == 0 {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
	}
	return string(runes)
}

const insertQuery = `
insert into users (given_name_encrypted, family_name_encrypted, email_encrypted, email_hash, inserted_at, updated_at, country, email_confirmed_at)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (email_hash) do nothing;
`

// This program decrypts the given age-encrypted CSV file
func main() {
	crypterer, err := crypto.NewEncryptionProvider(config.DatabaseEncryptionKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize database encryption provider: %s", err)
	}
	sqlcrypter.Init(crypterer)

	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var nationalityMapping map[string]string
	if err := json.Unmarshal(CountryMapping, &nationalityMapping); err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(os.Stdin)

	var rowIndex int

	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())

	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		// skip header
		rowIndex++
		if rowIndex == 1 {
			continue
		}

		givenName := capitalize(row[0])
		familyName := capitalize(row[1])
		email := strings.ToLower(row[2])
		emailVerified := row[5] == "Yes"
		insertedAt, err := time.Parse(time.RFC3339, row[9])
		if err != nil {
			log.Fatalf("Failed to parse time %v", row[9])
		}

		nationality := strings.TrimSpace(strings.ToLower(row[14]))
		iso, ok := nationalityMapping[nationality]
		if nationality != "" && !ok {
			log.Fatalf("Failed to guess nationality for %v", row[14])
		}

		var emailConfirmedAt *time.Time
		if emailVerified {
			emailConfirmedAt = &insertedAt
		}

		_, err = tx.Exec(
			context.Background(),
			insertQuery,
			sqlcrypter.NewEncryptedBytes(givenName),
			sqlcrypter.NewEncryptedBytes(familyName),
			sqlcrypter.NewEncryptedBytes(email),
			crypto.HashEmail(email),
			insertedAt,
			insertedAt,
			iso,
			emailConfirmedAt,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Fatal(err)
	}
}
