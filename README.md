# GoSalla - Salla Go SDK

A comprehensive Go SDK for the [Salla](https://salla.sa) e-commerce platform. This package provides easy-to-use interfaces for both API integration and webhook handling.

## Features

- ✅ **OAuth 2.0 Authentication** - Full OAuth flow with automatic token refresh
- ✅ **Complete API Coverage** - Products, Orders, Customers, Categories, and Brands
- ✅ **Webhook Support** - HMAC signature verification and typed event handlers
- ✅ **Pagination** - Built-in pagination support for list endpoints
- ✅ **Type-Safe** - Fully typed request and response structures
- ✅ **Zero External Dependencies** - Uses only Go standard library
- ✅ **Thread-Safe** - Safe for concurrent use

## Installation

```bash
go get github.com/abdalgaderserag/gosalla
```

## Quick Start

### 1. OAuth Authentication

```go
package main

import (
    "github.com/abdalgaderserag/gosalla"
    "log"
)

func main() {
    // Create OAuth config
    oauthConfig := &gosalla.OAuthConfig{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
        RedirectURI:  "your_redirect_uri",
        Scopes:       []string{"offline_access"},
    }
    
    // Generate authorization URL
    authURL := oauthConfig.GetAuthorizationURL("state")
    // Redirect user to authURL
    
    // After user authorizes, exchange code for token
    token, err := oauthConfig.ExchangeCode("authorization_code")
    if err != nil {
        log.Fatal(err)
    }
    
    // Token will be automatically refreshed when needed
}
```

### 2. Using the API Client

```go
package main

import (
    "github.com/abdalgaderserag/gosalla"
    "log"
)

func main() {
    // Create client with OAuth config and token
    client := gosalla.NewClient(oauthConfig, token)
    
    // List products
    products, pagination, err := client.Products.List(&gosalla.ListOptions{
        Page:    1,
        PerPage: 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a product
    newProduct := &gosalla.CreateProductRequest{
        Name:     "My Product",
        Price:    99.99,
        Quantity: 100,
        SKU:      "PROD-001",
    }
    
    product, err := client.Products.Create(newProduct)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Webhook Handling

```go
package main

import (
    "net/http"
    "github.com/abdalgaderserag/gosalla"
)

func main() {
    // Create webhook handler with your secret
    handler := gosalla.NewWebhookHandler("your_webhook_secret")
    
    // Register event handlers
    handler.OnProductCreated(func(event *gosalla.ProductWebhookEvent) error {
        // Handle product creation
        println("Product created:", event.Data.Name)
        return nil
    })
    
    handler.OnOrderCreated(func(event *gosalla.OrderWebhookEvent) error {
        // Handle order creation
        println("Order created:", event.Data.ReferenceID)
        return nil
    })
    
    // Start HTTP server
    http.Handle("/webhook", handler)
    http.ListenAndServe(":8080", nil)
}
```

## API Reference

### Client

```go
client := gosalla.NewClient(oauthConfig, token)
```

#### Products

```go
// List products
products, pagination, err := client.Products.List(opts)

// Get product by ID
product, err := client.Products.Get(id)

// Get product by SKU
product, err := client.Products.GetBySKU(sku)

// Create product
product, err := client.Products.Create(request)

// Update product
product, err := client.Products.Update(id, request)

// Delete product
err := client.Products.Delete(id)

// Change product status
err := client.Products.ChangeStatus(id, "active")
```

#### Orders

```go
// List orders
orders, pagination, err := client.Orders.List(opts)

// Get order by ID
order, err := client.Orders.Get(id)

// List order reservations
reservations, pagination, err := client.Orders.ListReservations(opts)
```

#### Customers

```go
// List customers
customers, pagination, err := client.Customers.List(opts)

// Get customer by ID
customer, err := client.Customers.Get(id)

// Create customer
customer, err := client.Customers.Create(request)

// Update customer
customer, err := client.Customers.Update(id, request)
```

#### Categories

```go
// List categories
categories, pagination, err := client.Categories.List(opts)

// Get category by ID
category, err := client.Categories.Get(id)

// Create category
category, err := client.Categories.Create(request)

// Update category
category, err := client.Categories.Update(id, request)

// Delete category
err := client.Categories.Delete(id)
```

#### Brands

```go
// List brands
brands, pagination, err := client.Brands.List(opts)

// Get brand by ID
brand, err := client.Brands.Get(id)

// Create brand
brand, err := client.Brands.Create(request)

// Update brand
brand, err := client.Brands.Update(id, request)

// Delete brand
err := client.Brands.Delete(id)
```

### Webhooks

#### Event Types

```go
const (
    EventProductCreated   = "product.created"
    EventProductUpdated   = "product.updated"
    EventProductDeleted   = "product.deleted"
    EventOrderCreated     = "order.created"
    EventOrderUpdated     = "order.updated"
    EventOrderCancelled   = "order.cancelled"
    EventCustomerCreated  = "customer.created"
    EventCustomerUpdated  = "customer.updated"
    // ... and more
)
```

#### Webhook Handler

```go
handler := gosalla.NewWebhookHandler(secret)

// Type-safe handlers
handler.OnProductCreated(func(event *ProductWebhookEvent) error {
    // event.Data is a Product struct
    return nil
})

handler.OnOrderCreated(func(event *OrderWebhookEvent) error {
    // event.Data is an Order struct
    return nil
})

handler.OnCustomerCreated(func(event *CustomerWebhookEvent) error {
    // event.Data is a Customer struct
    return nil
})

// Generic handler for any event type
handler.On("custom.event", func(event *WebhookEvent) error {
    // event.Data is map[string]interface{}
    return nil
})
```

## Examples

See the [`examples/`](./examples) directory for complete working examples:

- [`examples/oauth/`](./examples/oauth) - OAuth 2.0 authentication flow
- [`examples/products/`](./examples/products) - Product API operations
- [`examples/webhook/`](./examples/webhook) - Webhook server implementation

## Configuration

### Environment Variables

```bash
# OAuth Credentials
export SALLA_CLIENT_ID="your_client_id"
export SALLA_CLIENT_SECRET="your_client_secret"
export SALLA_REDIRECT_URI="https://yourdomain.com/callback"

# For API requests (after OAuth)
export SALLA_ACCESS_TOKEN="your_access_token"

# For webhooks
export SALLA_WEBHOOK_SECRET="your_webhook_secret"
```

## Error Handling

The SDK provides custom error types for better error handling:

```go
products, _, err := client.Products.List(nil)
if err != nil {
    if gosalla.IsNotFoundError(err) {
        // Handle 404
    } else if gosalla.IsUnauthorizedError(err) {
        // Handle 401 - maybe refresh token
    } else if gosalla.IsRateLimitError(err) {
        // Handle 429 - rate limited
    } else {
        // Handle other errors
    }
}
```

## Pagination

All list endpoints support pagination:

```go
opts := &gosalla.ListOptions{
    Page:    1,
    PerPage: 20,
}

products, pagination, err := client.Products.List(opts)
if err != nil {
    log.Fatal(err)
}

// Check if there are more pages
if pagination.HasNextPage() {
    nextPage := pagination.NextPage()
    // Fetch next page...
}
```

## Token Refresh

Tokens are automatically refreshed when needed:

```go
client := gosalla.NewClient(oauthConfig, token)

// The client will automatically refresh the token before it expires
// You can also manually refresh:
err := client.RefreshTokenIfNeeded()
if err != nil {
    // Handle error
}

// Get the current token (for persistence)
currentToken := client.GetToken()
```

## Testing

Run the tests:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -v -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Support

- **Documentation**: [Salla Developers](https://docs.salla.dev)
- **API Reference**: [Salla API Docs](https://docs.salla.dev/docs/merchant/)
- **Issues**: [GitHub Issues](https://github.com/abdalgaderserag/gosalla/issues)

## Acknowledgments

Built for the Salla e-commerce platform. Not officially affiliated with Salla.
