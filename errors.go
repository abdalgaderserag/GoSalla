package gosalla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the Salla API
type APIError struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Errors     map[string]interface{} `json:"errors,omitempty"`
	Response   *http.Response         `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("salla api error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("salla api error (status %d)", e.StatusCode)
}

// ErrorResponse represents the structure of error responses from Salla API
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// parseErrorResponse attempts to parse an error response from the API
func parseErrorResponse(resp *http.Response) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Response:   resp,
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiErr.Message = "failed to read error response"
		return apiErr
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// If we can't parse the error response, just use the raw body
		apiErr.Message = string(body)
		return apiErr
	}

	apiErr.Message = errResp.Message
	if errResp.Data != nil {
		apiErr.Errors = errResp.Data
	}

	return apiErr
}

// IsNotFoundError checks if the error is a 404 Not Found error
func IsNotFoundError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorizedError checks if the error is a 401 Unauthorized error
func IsUnauthorizedError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsRateLimitError checks if the error is a 429 Rate Limit error
func IsRateLimitError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}
