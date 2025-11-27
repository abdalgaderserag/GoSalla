/*
Package gosalla provides a Go client library for the Salla e-commerce platform API.

# Overview

The gosalla package offers a comprehensive, type-safe interface for integrating
with Salla's REST API and handling webhooks. It includes OAuth 2.0 authentication
with automatic token refresh, full CRUD operations for products, orders, customers,
categories, and brands, and webhook signature verification.

# Features

  - OAuth 2.0 authentication with automatic token refresh
  - Complete API coverage for core resources
  - Webhook handling with HMAC signature verification
  - Built-in pagination support
  - Thread-safe operations
  - Zero external dependencies

# Quick Start

OAuth Authentication:

	config := &gosalla.OAuthConfig{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		RedirectURI:  "your_redirect_uri",
	}

	authURL := config.GetAuthorizationURL("state")
	// Redirect user to authURL

	// After authorization
	token, err := config.ExchangeCode("authorization_code")
	if err != nil {
		log.Fatal(err)
	}

API Client Usage:

	client := gosalla.NewClient(config, token)

	// List products
	products, pagination, err := client.Products.List(&gosalla.ListOptions{
		Page:    1,
		PerPage: 10,
	})

	// Create a product
	product, err := client.Products.Create(&gosalla.CreateProductRequest{
		Name:     "My Product",
		Price:    99.99,
		Quantity: 100,
	})

Webhook Handling:

	handler := gosalla.NewWebhookHandler("webhook_secret")

	handler.OnProductCreated(func(event *gosalla.ProductWebhookEvent) error {
		fmt.Println("Product created:", event.Data.Name)
		return nil
	})

	http.Handle("/webhook", handler)
	http.ListenAndServe(":8080", nil)

# Error Handling

The package provides custom error types and helper functions:

	products, _, err := client.Products.List(nil)
	if err != nil {
		if gosalla.IsNotFoundError(err) {
			// Handle 404
		} else if gosalla.IsUnauthorizedError(err) {
			// Handle 401
		}
	}

# Resources

For more information, visit:
  - Salla Developer Documentation: https://docs.salla.dev
  - Salla API Reference: https://docs.salla.dev/docs/merchant/
*/
package gosalla
