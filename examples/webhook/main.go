package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	
	"github.com/abdalgaderserag/gosalla"
)

func main() {
	// Get webhook secret from environment variable
	webhookSecret := os.Getenv("SALLA_WEBHOOK_SECRET")
	
	if webhookSecret == "" {
		log.Println("Warning: SALLA_WEBHOOK_SECRET not set. Signature verification will be skipped.")
	}
	
	// Create webhook handler
	handler := gosalla.NewWebhookHandler(webhookSecret)
	
	// Register handlers for specific events
	handler.OnProductCreated(func(event *gosalla.ProductWebhookEvent) error {
		fmt.Printf("\n[Product Created] %s (ID: %d)\n", event.Data.Name, event.Data.ID)
		fmt.Printf("Price: %.2f, SKU: %s\n", event.Data.Price, event.Data.SKU)
		
		// Handle the product creation event
		// For example, sync with your inventory system
		
		return nil
	})
	
	handler.OnOrderCreated(func(event *gosalla.OrderWebhookEvent) error {
		fmt.Printf("\n[Order Created] Order #%s\n", event.Data.ReferenceID)
		fmt.Printf("Customer: %s (%s)\n", event.Data.Customer.Name, event.Data.Customer.Email)
		fmt.Printf("Total: %.2f %s\n", 
			event.Data.Amount.Total, 
			event.Data.Amount.CurrencyCode)
		fmt.Printf("Items: %d\n", len(event.Data.Items))
		
		// Handle the order creation event
		// For example, send confirmation email, update inventory
		
		return nil
	})
	
	handler.OnCustomerCreated(func(event *gosalla.CustomerWebhookEvent) error {
		fmt.Printf("\n[Customer Created] %s %s\n", 
			event.Data.FirstName, 
			event.Data.LastName)
		fmt.Printf("Email: %s, Phone: %s\n", 
			event.Data.Email, 
			event.Data.Phone)
		
		// Handle the customer creation event
		// For example, add to mailing list, send welcome email
		
		return nil
	})
	
	// Register a generic handler for other event types
	handler.On(gosalla.EventProductUpdated, func(event *gosalla.WebhookEvent) error {
		fmt.Printf("\n[Product Updated] Merchant: %d\n", event.Merchant)
		fmt.Printf("Data: %+v\n", event.Data)
		return nil
	})
	
	handler.On(gosalla.EventOrderShipped, func(event *gosalla.WebhookEvent) error {
		fmt.Printf("\n[Order Shipped] Merchant: %d\n", event.Merchant)
		fmt.Printf("Data: %+v\n", event.Data)
		return nil
	})
	
	// Set up HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.Handle("/webhook", handler)
	
	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Add a root endpoint with instructions
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>Salla Webhook Server</h1>
				<p>Webhook endpoint: <code>POST /webhook</code></p>
				<p>Health check: <code>GET /health</code></p>
				
				<h2>Test Webhook</h2>
				<p>You can test the webhook endpoint using curl:</p>
				<pre>
curl -X POST http://localhost:%s/webhook \
  -H "Content-Type: application/json" \
  -H "X-Signature: your_signature_here" \
  -d '{
    "event": "product.created",
    "merchant": 12345,
    "data": {
      "id": 1,
      "name": "Test Product",
      "price": 99.99,
      "sku": "TEST-001"
    },
    "created_at": "2024-01-01T00:00:00Z"
  }'
				</pre>
			</body>
			</html>
		`, port)
	})
	
	fmt.Printf("Starting webhook server on port %s...\n", port)
	fmt.Printf("Webhook endpoint: http://localhost:%s/webhook\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
