package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/stripe/stripe-go/v84"
)

type stripeWebhookController struct {
	orderService  *services.OrderService
	stripeService services.StripeService
}

func StripeWebhookController(db queries.DBTX, stripeClient services.StripeService) *stripeWebhookController {

	return &stripeWebhookController{
		stripeService: stripeClient,
		orderService:  services.NewOrderService(db, stripeClient),
	}
}

func (cc *stripeWebhookController) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := cc.stripeService.DecodeWebhook(body, r.Header.Get("Stripe-Signature"))
	if err != nil {
		log.Printf("Processing Stripe webhook failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event := event.(type) {
	case *stripe.CheckoutSession:
		_, err := cc.orderService.MarkOrderPaidByCheckoutSessionID(r.Context(), event.ID)
		if err != nil && errors.Is(err, services.ErrOrderAlreadyPaid) {
			log.Printf("Order already paid: %s", event.ID)
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
