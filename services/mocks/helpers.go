package mocks

import (
	"crypto/rand"
	"encoding/hex"

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
