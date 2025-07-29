package api

import (
	"net/http"
	"blessing_addon/handlers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/create-addon-session":
		handlers.CreateAddOnSessionHandler(w, r)
	case "/api/webhook-stripe":
		handlers.StripeWebhookHandler(w, r)
	case "/api/addon-sales-status":
		handlers.GetAddOnSalesStatusHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

