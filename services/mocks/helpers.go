package mocks

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/stripe/stripe-go/v84"
)

func GenerateCheckoutSession() *stripe.CheckoutSession {
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

func UniqueUser(db queries.DBTX, ctx context.Context, overrides ...func(params *types.SeedUserParams)) (*queries.User, error) {
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
