package gosalla

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Webhook event types
const (
	// Product events
	EventProductCreated = "product.created"
	EventProductUpdated = "product.updated"
	EventProductDeleted = "product.deleted"
	
	// Order events
	EventOrderCreated   = "order.created"
	EventOrderUpdated   = "order.updated"
	EventOrderCancelled = "order.cancelled"
	EventOrderShipped   = "order.shipped"
	EventOrderDelivered = "order.delivered"
	
	// Customer events
	EventCustomerCreated = "customer.created"
	EventCustomerUpdated = "customer.updated"
	EventCustomerDeleted = "customer.deleted"
	
	// Category events
	EventCategoryCreated = "category.created"
	EventCategoryUpdated = "category.updated"
	EventCategoryDeleted = "category.deleted"
	
	// Brand events
	EventBrandCreated = "brand.created"
	EventBrandUpdated = "brand.updated"
	EventBrandDeleted = "brand.deleted"
	
	// Cart events
	EventCartAbandoned = "cart.abandoned"
	EventCartRestored  = "cart.restored"
	
	// Payment events
	EventPaymentCompleted = "payment.completed"
	EventPaymentFailed    = "payment.failed"
	
	// Shipping events
	EventShipmentCreated = "shipment.created"
	EventShipmentUpdated = "shipment.updated"
)

// WebhookEvent represents a webhook event from Salla
type WebhookEvent struct {
	Event     string                 `json:"event"`
	Merchant  int                    `json:"merchant"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// ProductWebhookEvent represents a product-related webhook event
type ProductWebhookEvent struct {
	Event     string    `json:"event"`
	Merchant  int       `json:"merchant"`
	Data      Product   `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderWebhookEvent represents an order-related webhook event
type OrderWebhookEvent struct {
	Event     string    `json:"event"`
	Merchant  int       `json:"merchant"`
	Data      Order     `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

// CustomerWebhookEvent represents a customer-related webhook event
type CustomerWebhookEvent struct {
	Event     string    `json:"event"`
	Merchant  int       `json:"merchant"`
	Data      Customer  `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

// VerifyWebhookSignature verifies the HMAC signature of a webhook request
// The signature is typically sent in the X-Signature header
func VerifyWebhookSignature(secret string, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// ParseWebhook parses a webhook payload into a WebhookEvent
func ParseWebhook(payload []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &event, nil
}

// ParseProductWebhook parses a product webhook payload
func ParseProductWebhook(payload []byte) (*ProductWebhookEvent, error) {
	var event ProductWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse product webhook: %w", err)
	}
	return &event, nil
}

// ParseOrderWebhook parses an order webhook payload
func ParseOrderWebhook(payload []byte) (*OrderWebhookEvent, error) {
	var event OrderWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse order webhook: %w", err)
	}
	return &event, nil
}

// ParseCustomerWebhook parses a customer webhook payload
func ParseCustomerWebhook(payload []byte) (*CustomerWebhookEvent, error) {
	var event CustomerWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse customer webhook: %w", err)
	}
	return &event, nil
}

// WebhookHandler defines a function that handles webhook events
type WebhookHandler func(*WebhookEvent) error

// WebhookHandlerFunc is an HTTP handler function for processing webhooks
type WebhookHandlerFunc struct {
	Secret   string
	Handlers map[string]WebhookHandler
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(secret string) *WebhookHandlerFunc {
	return &WebhookHandlerFunc{
		Secret:   secret,
		Handlers: make(map[string]WebhookHandler),
	}
}

// On registers a handler for a specific event type
func (h *WebhookHandlerFunc) On(eventType string, handler WebhookHandler) {
	h.Handlers[eventType] = handler
}

// OnProductCreated registers a handler for product.created events
func (h *WebhookHandlerFunc) OnProductCreated(handler func(*ProductWebhookEvent) error) {
	h.On(EventProductCreated, func(event *WebhookEvent) error {
		productEvent, err := convertToProductEvent(event)
		if err != nil {
			return err
		}
		return handler(productEvent)
	})
}

// OnOrderCreated registers a handler for order.created events
func (h *WebhookHandlerFunc) OnOrderCreated(handler func(*OrderWebhookEvent) error) {
	h.On(EventOrderCreated, func(event *WebhookEvent) error {
		orderEvent, err := convertToOrderEvent(event)
		if err != nil {
			return err
		}
		return handler(orderEvent)
	})
}

// OnCustomerCreated registers a handler for customer.created events
func (h *WebhookHandlerFunc) OnCustomerCreated(handler func(*CustomerWebhookEvent) error) {
	h.On(EventCustomerCreated, func(event *WebhookEvent) error {
		customerEvent, err := convertToCustomerEvent(event)
		if err != nil {
			return err
		}
		return handler(customerEvent)
	})
}

// ServeHTTP implements http.Handler
func (h *WebhookHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	// Verify signature if secret is provided
	if h.Secret != "" {
		signature := r.Header.Get("X-Signature")
		if signature == "" {
			// Also check for Authorization header
			signature = r.Header.Get("Authorization")
		}
		
		if !VerifyWebhookSignature(h.Secret, body, signature) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}
	
	// Parse the webhook event
	event, err := ParseWebhook(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse webhook: %v", err), http.StatusBadRequest)
		return
	}
	
	// Find and execute the handler for this event type
	handler, exists := h.Handlers[event.Event]
	if !exists {
		// No handler registered for this event type, but still accept it
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Execute the handler
	if err := handler(event); err != nil {
		http.Error(w, fmt.Sprintf("Handler error: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

// Helper functions to convert generic events to typed events
func convertToProductEvent(event *WebhookEvent) (*ProductWebhookEvent, error) {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	
	var product Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}
	
	return &ProductWebhookEvent{
		Event:     event.Event,
		Merchant:  event.Merchant,
		Data:      product,
		CreatedAt: event.CreatedAt,
	}, nil
}

func convertToOrderEvent(event *WebhookEvent) (*OrderWebhookEvent, error) {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	
	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	
	return &OrderWebhookEvent{
		Event:     event.Event,
		Merchant:  event.Merchant,
		Data:      order,
		CreatedAt: event.CreatedAt,
	}, nil
}

func convertToCustomerEvent(event *WebhookEvent) (*CustomerWebhookEvent, error) {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	
	var customer Customer
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	
	return &CustomerWebhookEvent{
		Event:     event.Event,
		Merchant:  event.Merchant,
		Data:      customer,
		CreatedAt: event.CreatedAt,
	}, nil
}
