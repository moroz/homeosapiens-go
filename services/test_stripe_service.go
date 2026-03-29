package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/moroz/homeosapiens-go/types"
	"github.com/stripe/stripe-go/v84"
)

type mockStripeService struct{}

func generateRandomCheckoutSessionID() string {
	bytes := make([]byte, 12)
	_, _ = rand.Read(bytes)
	return "cs_test_" + hex.EncodeToString(bytes)
}

func (m mockStripeService) CreateCheckoutSession(ctx context.Context, dto *types.OrderDTO) (*stripe.CheckoutSession, error) {
	id := generateRandomCheckoutSessionID()

	return &stripe.CheckoutSession{
		ID:  id,
		URL: "https://checkout.stripe.com/c/pay/" + id,
	}, nil
}

func MockStripeService() StripeService {
	return &mockStripeService{}
}
