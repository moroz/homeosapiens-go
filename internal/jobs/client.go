package jobs

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func NewClient(db queries.DBTX) (*river.Client[pgx.Tx], error) {
	conn, ok := db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to cast database connection as *pgxpool.Pool, got: %T", db)
	}

	return river.NewClient(riverpgxv5.New(conn), &river.Config{})
}
