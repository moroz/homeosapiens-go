//go:generate mockery --name=StripeService --output=mocks --outpkg=mocks
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type StripeService interface {
	CreateCheckoutSession(context.Context, *types.OrderDTO) (*stripe.CheckoutSession, error)
	DecodeWebhook(payload []byte, signature string) (any, error)
}

type stripeService struct {
	client        *stripe.Client
	signingSecret string
}

func NewStripeService(secretKey, signingSecret string) StripeService {
	client := stripe.NewClient(secretKey)
	return &stripeService{
		client:        client,
		signingSecret: signingSecret,
	}
}

func (s *stripeService) CreateCheckoutSession(ctx context.Context, order *types.OrderDTO) (*stripe.CheckoutSession, error) {
	var lineItems []*stripe.CheckoutSessionCreateLineItemParams

	for _, item := range order.LineItems {
		lineItems = append(lineItems, &stripe.CheckoutSessionCreateLineItemParams{
			Quantity: stripe.Int64(int64(item.Quantity)),
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency:   stripe.String(strings.ToLower(order.Currency)),
				UnitAmount: stripe.Int64(item.EventPriceAmount.Mul(decimal.NewFromInt(100)).BigInt().Int64()),
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name: stripe.String(item.EventTitle),
				},
			},
		})
	}

	successUrl := config.PublicUrl + "/orders/success?session_id={CHECKOUT_SESSION_ID}"

	return s.client.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		ClientReferenceID: stripe.String(order.ID.String()),
		Mode:              stripe.String(stripe.CheckoutSessionModePayment),
		Currency:          &order.Currency,
		LineItems:         lineItems,
		UIMode:            stripe.String(stripe.CheckoutSessionUIModeHosted),
		CustomerEmail:     stripe.String(order.Email.String()),
		SuccessURL:        stripe.String(successUrl),
	})
}

func (s *stripeService) DecodeWebhook(payload []byte, signature string) (any, error) {
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		s.signingSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		return nil, fmt.Errorf("DecodeWebhook: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		var payload stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &payload); err != nil {
			return nil, fmt.Errorf("DecodeWebhook: %w", err)
		}
		return &payload, nil
	default:
		return nil, fmt.Errorf("DecodeWebhook: unsupported event type %s", event.Type)
	}
}
