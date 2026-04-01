package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStripeWebhook(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

}
