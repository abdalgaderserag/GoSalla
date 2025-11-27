package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/abdalgaderserag/gosalla"
)

func main() {
	// Get credentials from environment variables
	clientID := os.Getenv("SALLA_CLIENT_ID")
	clientSecret := os.Getenv("SALLA_CLIENT_SECRET")
	redirectURI := os.Getenv("SALLA_REDIRECT_URI")
	
	if clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Fatal("Please set SALLA_CLIENT_ID, SALLA_CLIENT_SECRET, and SALLA_REDIRECT_URI environment variables")
	}
	
	// Create OAuth config
	oauthConfig := &gosalla.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       []string{"offline_access"},
	}
	
	// Generate authorization URL
	state := "random_state_string" // In production, use a secure random string
	authURL := oauthConfig.GetAuthorizationURL(state)
	
	fmt.Println("Visit this URL to authorize the application:")
	fmt.Println(authURL)
	fmt.Println()
	
	// In a real application, you would:
	// 1. Redirect the user to the authorization URL
	// 2. Handle the callback to your redirect URI
	// 3. Extract the authorization code from the callback
	// 4. Exchange the code for an access token
	
	// For this example, we'll simulate receiving the authorization code
	fmt.Print("Enter the authorization code from the callback: ")
	var code string
	fmt.Scanln(&code)
	
	// Exchange the code for an access token
	token, err := oauthConfig.ExchangeCode(code)
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}
	
	fmt.Println("\nSuccessfully obtained access token!")
	fmt.Printf("Access Token: %s\n", token.AccessToken[:20]+"...")
	fmt.Printf("Token Type: %s\n", token.TokenType)
	fmt.Printf("Expires At: %s\n", token.Expiry)
	
	// Save the token for later use
	// In a real application, you would securely store this token
	
	// Example of refreshing the token
	if token.RefreshToken != "" {
		fmt.Println("\nRefreshing access token...")
		newToken, err := oauthConfig.RefreshToken(token.RefreshToken)
		if err != nil {
			log.Fatalf("Failed to refresh token: %v", err)
		}
		
		fmt.Println("Successfully refreshed access token!")
		fmt.Printf("New Access Token: %s\n", newToken.AccessToken[:20]+"...")
	}
}
