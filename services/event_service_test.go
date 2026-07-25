package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventService_UpdateEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
}
