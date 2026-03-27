package services

import (
	"context"
	"strings"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v84"
)

type StripeService interface {
	CreateCheckoutSession(context.Context, *types.OrderDTO) (*stripe.CheckoutSession, error)
}

type stripeService struct {
	client *stripe.Client
}

func NewStripeService(secretKey string) StripeService {
	client := stripe.NewClient(secretKey)
	return &stripeService{client}
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
