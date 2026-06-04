package mocks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
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
