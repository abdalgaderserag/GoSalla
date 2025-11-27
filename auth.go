package gosalla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// OAuth endpoints
	authorizationURL = "https://accounts.salla.sa/oauth2/auth"
	tokenURL         = "https://accounts.salla.sa/oauth2/token"
)

// OAuthConfig holds the OAuth 2.0 configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

// TokenResponse represents the response from the token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// Token represents an OAuth 2.0 token with expiry tracking
type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Expiry       time.Time
}

// Valid checks if the token is still valid (not expired)
func (t *Token) Valid() bool {
	return t.AccessToken != "" && time.Now().Before(t.Expiry)
}

// GetAuthorizationURL generates the OAuth authorization URL
func (c *OAuthConfig) GetAuthorizationURL(state string) string {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("redirect_uri", c.RedirectURI)
	params.Add("response_type", "code")
	params.Add("state", state)
	
	if len(c.Scopes) > 0 {
		scopes := ""
		for i, scope := range c.Scopes {
			if i > 0 {
				scopes += " "
			}
			scopes += scope
		}
		params.Add("scope", scopes)
	} else {
		params.Add("scope", "offline_access")
	}

	return fmt.Sprintf("%s?%s", authorizationURL, params.Encode())
}

// ExchangeCode exchanges an authorization code for an access token
func (c *OAuthConfig) ExchangeCode(code string) (*Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURI)
	data.Set("scope", "offline_access")

	return c.requestToken(data)
}

// RefreshToken refreshes an expired access token using the refresh token
func (c *OAuthConfig) RefreshToken(refreshToken string) (*Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("refresh_token", refreshToken)

	return c.requestToken(data)
}

// requestToken makes a request to the token endpoint
func (c *OAuthConfig) requestToken(data url.Values) (*Token, error) {
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	token := &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return token, nil
}
