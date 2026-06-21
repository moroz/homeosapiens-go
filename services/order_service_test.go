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

func TestOrderService_CreateOrder(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate orders, cart_line_items cascade")
	require.NoError(t, err)

	product, err := mocks.Product(db, ctx)
	require.NoError(t, err)

	event, err := mocks.Event(db, ctx, func(params *queries.UpsertEventParams) {
		params.ProductID = &product.ID
	})
	require.NoError(t, err)

	cartId := uuid.Must(uuid.NewV7())
	_, err = db.Exec(ctx, "insert into cart_line_items (cart_id, product_id) values ($1, $2)", cartId, product.ID)
	require.NoError(t, err)

	countBefore, err := count(db, ctx, "orders")
	assert.NoError(t, err)

	checkoutSession := mocks.CheckoutSession()

	stripe := mocks.NewMockStripeService(t)
	stripe.EXPECT().CreateCheckoutSession(mock.Anything, mock.Anything).Return(checkoutSession, nil)

	email := mocks.UniqueEmail()

	db.Exec(ctx, "truncate river_job")
	srv := services.NewOrderService(db, stripe)
	order, err := srv.CreateOrder(ctx, cartId, nil, &types.OrderParams{
		PreferredLocale:     "pl",
		Email:               email,
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

	countAfter, err := count(db, ctx, "orders")
	assert.NoError(t, err)
	assert.Equal(t, countBefore+1, countAfter)

	countAfter, err = count(db, ctx, "cart_line_items")
	assert.NoError(t, err)
	assert.Zero(t, countAfter)

	job := rivertest.RequireInserted(ctx, t, riverpgxv5.New(db), &jobs.SendOrderEmailArgs{}, nil)
	assert.Equal(t, jobs.OrderEmailTypeOrderConfirmation, job.Args.EmailType)

	// Creating an anonymous order creates a user with no password
	user, err := services.NewUserService(db).FindUserByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Empty(t, user.PasswordHash)

	// Placing an order does not grant access to products yet
	var hasAccess, eventRegistered bool
	err = db.QueryRow(
		ctx,
		`select exists (select from user_product_access where user_id = $1 and product_id = $2), exists (select from event_registrations where user_id = $1 and event_id = $3)`,
		user.ID, product.ID, event.ID,
	).Scan(&hasAccess, &eventRegistered)
	require.NoError(t, err)
	assert.False(t, hasAccess)
	assert.False(t, eventRegistered)
}

func TestOrderService_MarkOrderPaidByCheckoutSessionID(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	defer db.Close()

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	product, err := mocks.Product(db, ctx)
	require.NoError(t, err)

	event, err := mocks.Event(db, t.Context(), func(params *queries.UpsertEventParams) {
		params.ProductID = &product.ID
	})
	require.NoError(t, err)

	cartId := uuid.Must(uuid.NewV7())
	_, err = db.Exec(t.Context(), "insert into cart_line_items (cart_id, product_id) values ($1, $2)", cartId, product.ID)
	require.NoError(t, err)

	checkoutSession := mocks.CheckoutSession()
	stripe := mocks.NewMockStripeService(t)
	stripe.EXPECT().CreateCheckoutSession(mock.Anything, mock.Anything).Return(checkoutSession, nil)

	srv := services.NewOrderService(db, stripe)

	order, err := mocks.Order(&mocks.OrderInput{
		DB:      db,
		Context: ctx,
		Stripe:  stripe,
		CartID:  cartId,
		User:    user,
	})
	require.NoError(t, err)
	require.Equal(t, user.ID, order.UserID)
	require.Nil(t, order.PaidAt)

	_, err = db.Exec(ctx, "truncate river_job")

	updated, err := srv.MarkOrderPaidByCheckoutSessionID(ctx, checkoutSession.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updated.PaidAt)

	job := rivertest.RequireInserted(ctx, t, riverpgxv5.New(db), &jobs.SendOrderEmailArgs{}, nil)
	assert.Equal(t, jobs.OrderEmailTypePaymentConfirmation, job.Args.EmailType)

	var hasAccess, eventRegistered bool
	err = db.QueryRow(
		ctx,
		`select exists (select from user_product_access where user_id = $1 and product_id = $2), exists (select from event_registrations where user_id = $1 and event_id = $3)`,
		user.ID, product.ID, event.ID,
	).Scan(&hasAccess, &eventRegistered)
	require.NoError(t, err)
	assert.True(t, hasAccess)
	assert.True(t, eventRegistered)
}
