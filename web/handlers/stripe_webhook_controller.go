package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type stripeWebhookController struct {
	orderService *services.OrderService
}

func StripeWebhookController(db queries.DBTX, stripeClient services.StripeService) *stripeWebhookController {
	return &stripeWebhookController{
		orderService: services.NewOrderService(db, stripeClient),
	}
}

func (cc *stripeWebhookController) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		payload,
		r.Header.Get("Stripe-Signature"),
		config.StripeWebhookSigningSecret,
		// TODO: Can we get Stripe to actually support the latest version of their own API with their own Go SDK?
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		log.Printf("Webhook signature validation failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var payload stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &payload); err != nil {
			log.Printf("Error parsing webhook JSON: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = cc.orderService.MarkOrderPaidByCheckoutSessionID(r.Context(), payload.ID)
		if err != nil && errors.Is(err, services.ErrOrderAlreadyPaid) {
			log.Printf("Order already paid: %s", payload.ID)
			w.WriteHeader(http.StatusOK)
			return
		}
		if err != nil {
			log.Printf("Error marking order as paid: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
