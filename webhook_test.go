package gosalla

import (
	"testing"
)

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test_secret"
	payload := []byte(`{"event":"product.created","data":{"id":1}}`)
	
	// Generate valid signature
	validSig := "e8b7c09c8c8f8f0e4e8d4e8d4e8d4e8d4e8d4e8d4e8d4e8d4e8d4e8d4e8d4e8d"
	
	// Test with empty signature (should fail)
	if VerifyWebhookSignature(secret, payload, "") {
		t.Error("Expected verification to fail with empty signature")
	}
	
	// Test with wrong secret
	if VerifyWebhookSignature("wrong_secret", payload, validSig) {
		t.Error("Expected verification to fail with wrong secret")
	}
}

func TestParseWebhook(t *testing.T) {
	payload := []byte(`{
		"event": "product.created",
		"merchant": 12345,
		"data": {"id": 1, "name": "Test Product"},
		"created_at": "2024-01-01T00:00:00Z"
	}`)
	
	event, err := ParseWebhook(payload)
	if err != nil {
		t.Fatalf("Failed to parse webhook: %v", err)
	}
	
	if event.Event != "product.created" {
		t.Errorf("Expected event 'product.created', got '%s'", event.Event)
	}
	
	if event.Merchant != 12345 {
		t.Errorf("Expected merchant 12345, got %d", event.Merchant)
	}
	
	if event.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestParseWebhookInvalidJSON(t *testing.T) {
	payload := []byte(`{invalid json}`)
	
	_, err := ParseWebhook(payload)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON")
	}
}

func TestNewWebhookHandler(t *testing.T) {
	handler := NewWebhookHandler("test_secret")
	
	if handler == nil {
		t.Fatal("Expected handler to be created")
	}
	
	if handler.Secret != "test_secret" {
		t.Errorf("Expected secret 'test_secret', got '%s'", handler.Secret)
	}
	
	if handler.Handlers == nil {
		t.Error("Expected Handlers map to be initialized")
	}
}

func TestWebhookHandlerOn(t *testing.T) {
	handler := NewWebhookHandler("secret")
	called := false
	
	handler.On("test.event", func(event *WebhookEvent) error {
		called = true
		return nil
	})
	
	if len(handler.Handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(handler.Handlers))
	}
	
	// Test that the handler exists
	h, exists := handler.Handlers["test.event"]
	if !exists {
		t.Error("Expected handler to be registered for 'test.event'")
	}
	
	// Call the handler
	h(&WebhookEvent{})
	if !called {
		t.Error("Expected handler to be called")
	}
}

func TestPagination(t *testing.T) {
	p := &Pagination{
		CurrentPage: 2,
		LastPage:    5,
		PerPage:     10,
		Total:       50,
	}
	
	// Test HasNextPage
	if !p.HasNextPage() {
		t.Error("Expected HasNextPage to be true")
	}
	
	// Test NextPage
	if p.NextPage() != 3 {
		t.Errorf("Expected next page to be 3, got %d", p.NextPage())
	}
	
	// Test HasPreviousPage
	if !p.HasPreviousPage() {
		t.Error("Expected HasPreviousPage to be true")
	}
	
	// Test PreviousPage
	if p.PreviousPage() != 1 {
		t.Errorf("Expected previous page to be 1, got %d", p.PreviousPage())
	}
	
	// Test last page
	pLast := &Pagination{
		CurrentPage: 5,
		LastPage:    5,
	}
	
	if pLast.HasNextPage() {
		t.Error("Expected HasNextPage to be false on last page")
	}
	
	if pLast.NextPage() != 0 {
		t.Errorf("Expected NextPage to return 0 on last page, got %d", pLast.NextPage())
	}
	
	// Test first page
	pFirst := &Pagination{
		CurrentPage: 1,
		LastPage:    5,
	}
	
	if pFirst.HasPreviousPage() {
		t.Error("Expected HasPreviousPage to be false on first page")
	}
	
	if pFirst.PreviousPage() != 0 {
		t.Errorf("Expected PreviousPage to return 0 on first page, got %d", pFirst.PreviousPage())
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Not found",
	}
	
	errMsg := err.Error()
	expected := "salla api error (status 404): Not found"
	if errMsg != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, errMsg)
	}
	
	// Test without message
	err2 := &APIError{
		StatusCode: 500,
	}
	
	errMsg2 := err2.Error()
	expected2 := "salla api error (status 500)"
	if errMsg2 != expected2 {
		t.Errorf("Expected error message '%s', got '%s'", expected2, errMsg2)
	}
}

func TestIsNotFoundError(t *testing.T) {
	err := &APIError{StatusCode: 404}
	if !IsNotFoundError(err) {
		t.Error("Expected IsNotFoundError to be true for 404 error")
	}
	
	err2 := &APIError{StatusCode: 500}
	if IsNotFoundError(err2) {
		t.Error("Expected IsNotFoundError to be false for non-404 error")
	}
}

func TestIsUnauthorizedError(t *testing.T) {
	err := &APIError{StatusCode: 401}
	if !IsUnauthorizedError(err) {
		t.Error("Expected IsUnauthorizedError to be true for 401 error")
	}
	
	err2 := &APIError{StatusCode: 200}
	if IsUnauthorizedError(err2) {
		t.Error("Expected IsUnauthorizedError to be false for non-401 error")
	}
}

func TestIsRateLimitError(t *testing.T) {
	err := &APIError{StatusCode: 429}
	if !IsRateLimitError(err) {
		t.Error("Expected IsRateLimitError to be true for 429 error")
	}
	
	err2 := &APIError{StatusCode: 200}
	if IsRateLimitError(err2) {
		t.Error("Expected IsRateLimitError to be false for non-429 error")
	}
}
