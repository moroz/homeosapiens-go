package mocks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v84"
)

func CheckoutSession() *stripe.CheckoutSession {
	serial := make([]byte, 12)
	_, _ = rand.Read(serial)
	id := "cs_test_" + hex.EncodeToString(serial)
	return &stripe.CheckoutSession{
		ID:  id,
		URL: "https://checkout.stripe.com/c/pay/" + id,
	}
}

func UniqueEmail() string {
	unique := make([]byte, 4)
	_, _ = rand.Read(unique)
	return "user-" + hex.EncodeToString(unique) + "@example.com"
}

func User(db queries.DBTX, ctx context.Context, overrides ...func(params *types.SeedUserParams)) (*queries.User, error) {
	email := UniqueEmail()
	params := &types.SeedUserParams{
		Email:      email,
		GivenName:  "John",
		FamilyName: "Smith",
	}

	for _, f := range overrides {
		f(params)
	}

	return services.NewUserService(db).CreateUser(ctx, params)
}

func Product(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.InsertProductParams)) (*queries.Product, error) {
	params := &queries.InsertProductParams{
		ProductType:       queries.ProductTypeEvent,
		TitlePl:           "Some product",
		TitleEn:           "Some product",
		BasePriceAmount:   decimal.NewFromInt(0),
		BasePriceCurrency: "EUR",
	}

	for _, f := range overrides {
		f(params)
	}

	return queries.New(db).InsertProduct(ctx, params)
}

func Event(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.UpsertEventParams)) (*queries.Event, error) {
	unique := make([]byte, 4)
	_, _ = rand.Read(unique)

	params := &queries.UpsertEventParams{
		ID:            uuid.Must(uuid.NewV7()),
		EventType:     queries.EventTypeSeminar,
		TitleEn:       "Some event",
		TitlePl:       "Some event",
		Slug:          "event-" + hex.EncodeToString(unique),
		StartsAt:      time.Now(),
		EndsAt:        time.Now().Add(2 * time.Hour),
		IsVirtual:     true,
		DescriptionEn: "Some description",
	}

	for _, f := range overrides {
		f(params)
	}

	return queries.New(db).UpsertEvent(ctx, params)
}

func PaidEvent(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.UpsertEventParams)) (*queries.Event, error) {
	product, err := Product(db, ctx, func(p *queries.InsertProductParams) {
		p.BasePriceAmount = decimal.NewFromInt(560)
	})
	if err != nil {
		return nil, err
	}
	combined := append([]func(*queries.UpsertEventParams){
		func(p *queries.UpsertEventParams) {
			p.ProductID = &product.ID
		},
	}, overrides...)
	return Event(db, ctx, combined...)
}

func VideoGroup(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.InsertVideoGroupParams)) (*queries.VideoGroup, error) {
	unique := make([]byte, 4)
	_, _ = rand.Read(unique)

	params := &queries.InsertVideoGroupParams{
		TitleEn: "Some video group",
		TitlePl: "Some video group",
		Slug:    "video-group-" + hex.EncodeToString(unique),
	}

	for _, f := range overrides {
		f(params)
	}

	return queries.New(db).InsertVideoGroup(ctx, params)
}

func Video(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.InsertVideoParams)) (*queries.Video, error) {
	unique := make([]byte, 4)
	_, _ = rand.Read(unique)

	params := &queries.InsertVideoParams{
		TitlePl:         "Some video title",
		TitleEn:         "Some video title",
		Slug:            "video-" + hex.EncodeToString(unique),
		Provider:        queries.VideoProviderCloudfront,
		DurationSeconds: new(int32(420)),
		RecordedOn: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	}

	for _, f := range overrides {
		f(params)
	}

	return queries.New(db).InsertVideo(ctx, params)
}

func EventRegistration(db queries.DBTX, ctx context.Context, event *queries.Event, user *queries.User) (*queries.EventRegistration, error) {
	return queries.New(db).InsertEventRegistration(ctx, &queries.InsertEventRegistrationParams{
		EventID: event.ID,
		UserID:  user.ID,
	})
}

// Cart creates cart line items for a fresh cart and returns its id. Pass the
// product ids the cart should contain (one line item each, quantity 1).
func Cart(db queries.DBTX, ctx context.Context, productIDs ...uuid.UUID) (uuid.UUID, error) {
	cartID := uuid.Must(uuid.NewV7())
	for _, pid := range productIDs {
		_, err := db.Exec(ctx, "insert into cart_line_items (cart_id, product_id) values ($1, $2)", cartID, pid)
		if err != nil {
			return cartID, err
		}
	}
	return cartID, nil
}

type OrderInput struct {
	DB      *pgxpool.Pool
	Context context.Context
	Stripe  services.StripeService
	CartID  uuid.UUID
	User    *queries.User
	// Overrides mutate the default Polish billing params before the order is
	// created, e.g. to pin the buyer's email.
	Overrides []func(params *types.OrderParams)
}

// Order creates a paid-order draft from props.CartID by running it through
// OrderService.CreateOrder, with sane Polish billing defaults.
func Order(props *OrderInput) (*types.OrderDTO, error) {
	params := &types.OrderParams{
		PreferredLocale:     "pl",
		Email:               UniqueEmail(),
		BillingGivenName:    "John",
		BillingFamilyName:   "Smith",
		BillingPhone:        "+48555123456",
		BillingAddressLine1: "ul. Półwiejska 20",
		BillingCity:         "Poznań",
		BillingPostalCode:   "12-345",
		BillingCountry:      "PL",
	}

	for _, f := range props.Overrides {
		f(params)
	}

	return services.NewOrderService(props.DB, props.Stripe).CreateOrder(props.Context, props.CartID, props.User, params)
}

func UserSession(db queries.DBTX, ctx context.Context, user *queries.User) (sessions.Payload, error) {
	token, err := services.NewUserTokenService(db).IssueAccessTokenForUser(ctx, user, config.AccessTokenValidity)
	if err != nil {
		return nil, err
	}
	return sessions.Payload{
		config.AccessTokenSessionKey: token.Token,
		config.LanguageSessionKey:    "en",
	}, nil
}
