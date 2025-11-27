package gosalla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Salla API
	DefaultBaseURL = "https://api.salla.dev/admin/v2"
	
	// DefaultUserAgent is the default user agent for requests
	DefaultUserAgent = "gosalla/1.0"
)

// Client is the main client for interacting with the Salla API
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
	
	// OAuth configuration and token
	oauthConfig *OAuthConfig
	token       *Token
	tokenMu     sync.RWMutex
	
	// API resource clients
	Products   *ProductsService
	Orders     *OrdersService
	Customers  *CustomersService
	Categories *CategoriesService
	Brands     *BrandsService
}

// NewClient creates a new Salla API client
func NewClient(oauthConfig *OAuthConfig, token *Token) *Client {
	c := &Client{
		baseURL:     DefaultBaseURL,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		userAgent:   DefaultUserAgent,
		oauthConfig: oauthConfig,
		token:       token,
	}
	
	// Initialize service clients
	c.Products = &ProductsService{client: c}
	c.Orders = &OrdersService{client: c}
	c.Customers = &CustomersService{client: c}
	c.Categories = &CategoriesService{client: c}
	c.Brands = &BrandsService{client: c}
	
	return c
}

// SetBaseURL sets a custom base URL for the API
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

// SetUserAgent sets a custom user agent
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// GetToken returns the current access token (thread-safe)
func (c *Client) GetToken() *Token {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token
}

// SetToken sets a new access token (thread-safe)
func (c *Client) SetToken(token *Token) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = token
}

// RefreshTokenIfNeeded refreshes the access token if it's expired or about to expire
func (c *Client) RefreshTokenIfNeeded() error {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	
	// Check if token is still valid (with 5-minute buffer)
	if c.token != nil && time.Now().Add(5*time.Minute).Before(c.token.Expiry) {
		return nil
	}
	
	if c.token == nil || c.token.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}
	
	newToken, err := c.oauthConfig.RefreshToken(c.token.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	
	c.token = newToken
	return nil
}

// newRequest creates a new HTTP request with proper headers and authentication
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}
	
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	
	// Add authorization header
	c.tokenMu.RLock()
	if c.token != nil && c.token.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
	}
	c.tokenMu.RUnlock()
	
	return req, nil
}

// do executes an HTTP request and handles the response
func (c *Client) do(req *http.Request, v interface{}) error {
	// Refresh token if needed before making the request
	if err := c.RefreshTokenIfNeeded(); err != nil {
		// Update authorization header with new token
		c.tokenMu.RLock()
		if c.token != nil && c.token.AccessToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))
		}
		c.tokenMu.RUnlock()
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp)
	}
	
	// Parse response if a destination is provided
	if v != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}
	
	return nil
}

// Response represents a standard API response wrapper
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
