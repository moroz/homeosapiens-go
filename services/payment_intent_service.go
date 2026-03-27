package services

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/stripe/stripe-go/v84"
)

type PaymentIntentService interface {
	CreatePaymentIntentForOrder(*queries.Order) (*stripe.PaymentIntent, error)
}

type paymentIntentService struct {
}

func NewPaymentIntentService() PaymentIntentService {
	return &paymentIntentService{}
}

func (s *paymentIntentService) CreatePaymentIntentForOrder(order *queries.Order) (*stripe.PaymentIntent, error) {
	panic("unimplemented!")
}
