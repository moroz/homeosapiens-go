package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/stretchr/testify/assert"
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
	_ = db

	eventId := uuid.MustParse("019c5c9a-c5a4-7518-8317-65ae90516726")
	cartId := uuid.Must(uuid.NewV7())

	_, err = queries.New(db).InsertCartLineItem(t.Context(), &queries.InsertCartLineItemParams{
		CartID:  cartId,
		EventID: eventId,
	})
	require.NoError(t, err)

	countBefore, err := count(db, t.Context(), "orders")
	assert.NoError(t, err)

	srv := services.NewOrderService(db)
	order, err := srv.CreateOrder(t.Context(), cartId, nil, &types.OrderParams{
		Email:                  "user@example.com",
		BillingGivenName:       "John",
		BillingFamilyName:      "Smith",
		BillingPhone:           new("+48555123456"),
		BillingStreet:          "ul. Półwiejska",
		BillingHouseNumber:     "20",
		BillingApartmentNumber: nil,
		BillingCity:            "Poznań",
		BillingPostalCode:      new("12-345"),
		BillingCountry:         "PL",
	})

	assert.NotNil(t, order)

	countAfter, err := count(db, t.Context(), "orders")
	assert.NoError(t, err)
	assert.Equal(t, countBefore+1, countAfter)
}
