package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
    "github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
	"blessing_addon/store"
)

type CreateSessionRequest struct {
	UserID          string `json:"userId"`
	AmountCents     int64  `json:"amount_cents"`
	Currency        string `json:"currency"`
	StripeAccountID string `json:"stripe_account_id"`
}

func CreateAddOnSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(req.AmountCents * 20 / 100),
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
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://your-success-url.com"),
		CancelURL:  stripe.String("https://your-cancel-url.com"),
	}

	// Set the Stripe account for the transfer
	if req.StripeAccountID != "" {
		params.PaymentIntentData.TransferData = &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
			Destination: stripe.String(req.StripeAccountID),
		}
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("Stripe session error: %v", err)
		http.Error(w, "Stripe session creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"sessionId": s.ID,
	})
}

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Webhook error: %v", err)
		http.Error(w, "Webhook verification failed", http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		var s stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
			http.Error(w, "Invalid session data", http.StatusBadRequest)
			return
		}

		store.RecordAddOnPurchase(int(s.AmountTotal))
		log.Println("Blessing PDF purchased")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Blessing PDF purchase recorded",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAddOnSalesStatusHandler(w http.ResponseWriter, r *http.Request) {
	vendorID := r.URL.Query().Get("vendorId")
	if vendorID == "" {
		http.Error(w, "vendorId required", http.StatusBadRequest)
		return
	}

	stats := store.GetStats()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"totalAddOnSales":           stats.TotalSales,
		"totalAddOnRevenueCents":    stats.TotalRevenueCents,
		"totalAddOnCommissionCents": stats.TotalCommissionCents,
	})
}
