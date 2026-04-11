package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func count(db queries.DBTX, ctx context.Context, table string) (int, error) {
	var val int
	err := db.QueryRow(ctx, "select count(*) from "+table).Scan(&val)
	return val, err
}

func TestCreateOrder(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(t.Context(), "truncate orders, cart_line_items cascade")

	eventId := uuid.MustParse("019c5c9a-c5a4-7518-8317-65ae90516726")
	cartId := uuid.Must(uuid.NewV7())

	_, err = queries.New(db).InsertCartLineItem(t.Context(), &queries.InsertCartLineItemParams{
		CartID:  cartId,
		EventID: eventId,
	})
	require.NoError(t, err)

	countBefore, err := count(db, t.Context(), "orders")
	assert.NoError(t, err)

	cs := mocks.GenerateCheckoutSession()

	stripe := mocks.NewMockStripeService(t)
	stripe.EXPECT().CreateCheckoutSession(mock.Anything, mock.Anything).Return(cs, nil)

	db.Exec(t.Context(), "truncate river_job")
	srv := services.NewOrderService(db, stripe)
	order, err := srv.CreateOrder(t.Context(), cartId, nil, &types.OrderParams{
		PreferredLocale:     "pl",
		Email:               "user@example.com",
		BillingGivenName:    "John",
		BillingFamilyName:   "Smith",
		BillingPhone:        "+48555123456",
		BillingAddressLine1: "ul. Półwiejska 20",
		BillingCity:         "Poznań",
		BillingPostalCode:   "12-345",
		BillingCountry:      "PL",
	})

	require.NoError(t, err)
	require.NotNil(t, order)

	countAfter, err := count(db, t.Context(), "orders")
	assert.NoError(t, err)
	assert.Equal(t, countBefore+1, countAfter)

	countAfter, err = count(db, t.Context(), "cart_line_items")
	assert.NoError(t, err)
	assert.Zero(t, countAfter)

	rivertest.RequireInserted(t.Context(), t, riverpgxv5.New(db), &jobs.SendOrderEmailArgs{}, nil)
}
