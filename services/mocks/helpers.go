package mocks

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
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

func UniqueUser(db queries.DBTX, ctx context.Context, overrides ...func(params *queries.UpsertUserFromSeedDataParams)) (*queries.User, error) {
	email := UniqueEmail()
	params := &queries.UpsertUserFromSeedDataParams{
		Email:      sqlcrypter.NewEncryptedBytes(email),
		EmailHash:  crypto.HashEmail(email),
		GivenName:  sqlcrypter.NewEncryptedBytes("John"),
		FamilyName: sqlcrypter.NewEncryptedBytes("Smith"),
	}

	for _, f := range overrides {
		f(params)
	}

	return queries.New(db).UpsertUserFromSeedData(ctx, params)
}
