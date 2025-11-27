package gosalla

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := &OAuthConfig{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}
	
	token := &Token{
		AccessToken: "test_token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}
	
	client := NewClient(config, token)
	
	if client == nil {
		t.Fatal("Expected client to be created")
	}
	
	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", DefaultBaseURL, client.baseURL)
	}
	
	if client.userAgent != DefaultUserAgent {
		t.Errorf("Expected user agent %s, got %s", DefaultUserAgent, client.userAgent)
	}
	
	if client.Products == nil {
		t.Error("Expected Products service to be initialized")
	}
	
	if client.Orders == nil {
		t.Error("Expected Orders service to be initialized")
	}
	
	if client.Customers == nil {
		t.Error("Expected Customers service to be initialized")
	}
	
	if client.Categories == nil {
		t.Error("Expected Categories service to be initialized")
	}
	
	if client.Brands == nil {
		t.Error("Expected Brands service to be initialized")
	}
}

func TestSetBaseURL(t *testing.T) {
	client := NewClient(&OAuthConfig{}, &Token{})
	customURL := "https://custom.api.url"
	
	client.SetBaseURL(customURL)
	
	if client.baseURL != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.baseURL)
	}
}

func TestSetUserAgent(t *testing.T) {
	client := NewClient(&OAuthConfig{}, &Token{})
	customAgent := "CustomAgent/1.0"
	
	client.SetUserAgent(customAgent)
	
	if client.userAgent != customAgent {
		t.Errorf("Expected user agent %s, got %s", customAgent, client.userAgent)
	}
}

func TestGetSetToken(t *testing.T) {
	client := NewClient(&OAuthConfig{}, &Token{AccessToken: "initial"})
	
	// Test Get
	token := client.GetToken()
	if token.AccessToken != "initial" {
		t.Errorf("Expected access token 'initial', got '%s'", token.AccessToken)
	}
	
	// Test Set
	newToken := &Token{AccessToken: "new_token"}
	client.SetToken(newToken)
	
	retrievedToken := client.GetToken()
	if retrievedToken.AccessToken != "new_token" {
		t.Errorf("Expected access token 'new_token', got '%s'", retrievedToken.AccessToken)
	}
}

func TestTokenValid(t *testing.T) {
	// Valid token
	validToken := &Token{
		AccessToken: "test",
		Expiry:      time.Now().Add(1 * time.Hour),
	}
	
	if !validToken.Valid() {
		t.Error("Expected token to be valid")
	}
	
	// Expired token
	expiredToken := &Token{
		AccessToken: "test",
		Expiry:      time.Now().Add(-1 * time.Hour),
	}
	
	if expiredToken.Valid() {
		t.Error("Expected token to be invalid")
	}
	
	// Empty token
	emptyToken := &Token{}
	
	if emptyToken.Valid() {
		t.Error("Expected empty token to be invalid")
	}
}

func TestNewRequest(t *testing.T) {
	client := NewClient(&OAuthConfig{}, &Token{AccessToken: "test_token"})
	
	req, err := client.newRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Check method
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}
	
	// Check URL
	expectedURL := DefaultBaseURL + "/test"
	if req.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
	}
	
	// Check headers
	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type header to be application/json")
	}
	
	if req.Header.Get("Accept") != "application/json" {
		t.Error("Expected Accept header to be application/json")
	}
	
	if req.Header.Get("User-Agent") != DefaultUserAgent {
		t.Errorf("Expected User-Agent %s, got %s", DefaultUserAgent, req.Header.Get("User-Agent"))
	}
	
	authHeader := req.Header.Get("Authorization")
	expectedAuth := "Bearer test_token"
	if authHeader != expectedAuth {
		t.Errorf("Expected Authorization header %s, got %s", expectedAuth, authHeader)
	}
}

func TestNewRequestWithBody(t *testing.T) {
	client := NewClient(&OAuthConfig{}, &Token{})
	
	body := map[string]string{"key": "value"}
	req, err := client.newRequest("POST", "/test", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if req.Body == nil {
		t.Error("Expected request to have a body")
	}
}

func TestDoWithError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success": false, "message": "Not found"}`))
	}))
	defer server.Close()
	
	client := NewClient(&OAuthConfig{}, &Token{AccessToken: "test"})
	client.SetBaseURL(server.URL)
	
	req, _ := client.newRequest("GET", "/test", nil)
	err := client.do(req, nil)
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if !IsNotFoundError(err) {
		t.Error("Expected NotFoundError")
	}
}
