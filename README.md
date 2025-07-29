# blessing_addon
A Go-based backend module for Stripe Add-On Purchases using Stripe Connect. Handles checkout sessions, webhook processing, and sales analytics for Blessing PDF add-ons. Built for PowerOfAum assessment

# Blessing Add-On Module

This project implements the **Module B: Add-On Purchase – Blessing PDF** for PowerOfAum's backend assessment.

### Tech Stack

- Language: **Golang**
- Billing: **Stripe Checkout (Connect - Custom Accounts)**
- Storage: In-memory mock store
- Deployment: **Vercel**

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/create-addon-session` | POST | Creates Stripe Checkout Session |
| `/api/webhook-stripe` | POST | Handles Stripe webhook events |
| `/api/addon-sales-status?vendorId=...` | GET | Returns sales summary |

### Test Card

Use `4242 4242 4242 4242` with any future expiry and CVV for test payments.

---

###
Ayushi Katheria  
Made with Golang

