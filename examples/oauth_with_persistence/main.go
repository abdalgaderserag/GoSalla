package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abdalgaderserag/gosalla"
	_ "github.com/mattn/go-sqlite3" // SQLite driver (you can replace with MySQL/PostgreSQL)
)

// TokenManager handles token storage and automatic refresh
type TokenManager struct {
	db          *sql.DB
	oauthConfig *gosalla.OAuthConfig
}

// NewTokenManager creates a new token manager with database connection
func NewTokenManager(db *sql.DB, oauthConfig *gosalla.OAuthConfig) *TokenManager {
	return &TokenManager{
		db:          db,
		oauthConfig: oauthConfig,
	}
}

// InitDatabase creates the tokens table if it doesn't exist
func (tm *TokenManager) InitDatabase() error {
	query := `
		CREATE TABLE IF NOT EXISTS tokens (
			user_id TEXT PRIMARY KEY,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			token_type TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := tm.db.Exec(query)
	return err
}

// SaveToken stores or updates a token in the database
func (tm *TokenManager) SaveToken(userID string, token *gosalla.Token) error {
	query := `
		INSERT INTO tokens (user_id, access_token, refresh_token, token_type, expires_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id) DO UPDATE SET
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			token_type = excluded.token_type,
			expires_at = excluded.expires_at,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := tm.db.Exec(query,
		userID,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.Expiry,
	)

	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	log.Printf("[TokenManager] Token saved for user: %s", userID)
	return nil
}

// GetToken retrieves a token from the database and automatically refreshes if expired
func (tm *TokenManager) GetToken(userID string) (*gosalla.Token, error) {
	query := `
		SELECT access_token, refresh_token, token_type, expires_at
		FROM tokens
		WHERE user_id = ?
	`

	var token gosalla.Token
	err := tm.db.QueryRow(query, userID).Scan(
		&token.AccessToken,
		&token.RefreshToken,
		&token.TokenType,
		&token.Expiry,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no token found for user: %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Check if token is expired or about to expire (within 5 minutes)
	if time.Now().Add(5 * time.Minute).After(token.Expiry) {
		log.Printf("[TokenManager] Token expired or expiring soon, refreshing...")
		
		newToken, err := tm.RefreshToken(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		
		return newToken, nil
	}

	log.Printf("[TokenManager] Token retrieved for user: %s (valid until %s)", userID, token.Expiry.Format(time.RFC3339))
	return &token, nil
}

// RefreshToken refreshes an expired token and saves it to the database
func (tm *TokenManager) RefreshToken(userID string) (*gosalla.Token, error) {
	// Get current token
	query := `
		SELECT refresh_token
		FROM tokens
		WHERE user_id = ?
	`

	var refreshToken string
	err := tm.db.QueryRow(query, userID).Scan(&refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Refresh using OAuth config
	newToken, err := tm.oauthConfig.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Save new token
	if err := tm.SaveToken(userID, newToken); err != nil {
		return nil, err
	}

	log.Printf("[TokenManager] Token refreshed successfully for user: %s", userID)
	return newToken, nil
}

// DeleteToken removes a token from the database
func (tm *TokenManager) DeleteToken(userID string) error {
	query := `DELETE FROM tokens WHERE user_id = ?`
	_, err := tm.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	log.Printf("[TokenManager] Token deleted for user: %s", userID)
	return nil
}

// GetOrCreateClient gets a Salla API client with automatic token refresh
func (tm *TokenManager) GetOrCreateClient(userID string) (*gosalla.Client, error) {
	token, err := tm.GetToken(userID)
	if err != nil {
		return nil, err
	}

	client := gosalla.NewClient(tm.oauthConfig, token)
	return client, nil
}

// Example usage
func main() {
	// Setup database connection
	db, err := sql.Open("sqlite3", "./salla_tokens.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Setup OAuth config
	oauthConfig := &gosalla.OAuthConfig{
		ClientID:     os.Getenv("SALLA_CLIENT_ID"),
		ClientSecret: os.Getenv("SALLA_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SALLA_REDIRECT_URI"),
		Scopes:       []string{"offline_access"},
	}

	// Create token manager
	tokenManager := NewTokenManager(db, oauthConfig)

	// Initialize database tables
	if err := tokenManager.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Simulate user ID (in a real app, this comes from your auth system)
	userID := "user123"

	// Example 1: First time - Get token via OAuth and save
	fmt.Println("=== Example 1: Initial OAuth Flow ===")
	
	// Check if user already has a token
	existingToken, err := tokenManager.GetToken(userID)
	if err != nil {
		fmt.Println("No existing token found. Starting OAuth flow...")
		
		// Generate authorization URL
		authURL := oauthConfig.GetAuthorizationURL("random_state")
		fmt.Println("Visit this URL to authorize:")
		fmt.Println(authURL)
		fmt.Println()
		
		// Get authorization code from user
		fmt.Print("Enter the authorization code: ")
		var code string
		fmt.Scanln(&code)
		
		// Exchange code for token
		initialToken, err := oauthConfig.ExchangeCode(code)
		if err != nil {
			log.Fatalf("Failed to exchange code: %v", err)
		}
		
		// Save to database
		if err := tokenManager.SaveToken(userID, initialToken); err != nil {
			log.Fatalf("Failed to save token: %v", err)
		}
		
		fmt.Println("✓ Token saved to database")
	} else {
		fmt.Printf("✓ Token already exists (expires: %s)\n", existingToken.Expiry.Format(time.RFC3339))
	}

	// Example 2: Get client with automatic token refresh
	fmt.Println("\n=== Example 2: Using Client with Auto-Refresh ===")
	
	client, err := tokenManager.GetOrCreateClient(userID)
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Use the client - token will auto-refresh if needed
	products, pagination, err := client.Products.List(&gosalla.ListOptions{
		Page:    1,
		PerPage: 5,
	})
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}

	fmt.Printf("✓ Fetched %d products (page %d of %d)\n",
		len(products), pagination.CurrentPage, pagination.LastPage)

	for i, product := range products {
		fmt.Printf("  %d. %s - %.2f SAR\n", i+1, product.Name, product.Price)
	}

	// Example 3: Manual refresh (usually not needed)
	fmt.Println("\n=== Example 3: Manual Token Refresh ===")
	
	refreshedToken, err := tokenManager.RefreshToken(userID)
	if err != nil {
		log.Printf("Failed to refresh: %v", err)
	} else {
		fmt.Printf("✓ Token manually refreshed (new expiry: %s)\n", 
			refreshedToken.Expiry.Format(time.RFC3339))
	}

	fmt.Println("\n=== All operations completed successfully ===")
}
