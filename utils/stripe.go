package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY") // Put your secret in .env
}

// Store mock add-ons in-memory
var addOns = []map[string]interface{}{}

type AddOnRequest struct {
	UserID           string `json:"userId"`
	AmountCents      int64  `json:"amount_cents"`
	Currency         string `json:"currency"`
	StripeAccountID  string `json:"stripe_account_id"`
}

func CreateAddOnSession(w http.ResponseWriter, r *http.Request) {
	var req AddOnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create session params with 20% commission (application_fee)
	applicationFee := req.AmountCents * 20 / 100

	params := &stripe.CheckoutSessionParams{
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(applicationFee),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(req.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Blessing PDF"),
					},
					UnitAmount: stripe.Int64(req.AmountCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://your-success-url.com"),
		CancelURL:  stripe.String("https://your-cancel-url.com"),
	}

	params.SetStripeAccount(req.StripeAccountID)

	s, err := session.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"success":   true,
		"sessionId": s.ID,
	}
	json.NewEncoder(w).Encode(resp)
}

func StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusServiceUnavailable)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		http.Error(w, "Webhook verification failed", http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		// Simulate Blessing PDF storage
		addOns = append(addOns, map[string]interface{}{
			"session": event.ID,
		})

		resp := map[string]interface{}{
			"success": true,
			"message": "Blessing PDF purchased",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}
