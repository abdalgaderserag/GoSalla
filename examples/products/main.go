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
	accessToken := os.Getenv("SALLA_ACCESS_TOKEN")
	
	if clientID == "" || clientSecret == "" || accessToken == "" {
		log.Fatal("Please set SALLA_CLIENT_ID, SALLA_CLIENT_SECRET, and SALLA_ACCESS_TOKEN environment variables")
	}
	
	// Create OAuth config
	oauthConfig := &gosalla.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	
	// Create a token (in a real app, you would load this from secure storage)
	token := &gosalla.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}
	
	// Create the Salla API client
	client := gosalla.NewClient(oauthConfig, token)
	
	// List all products
	fmt.Println("Listing products...")
	products, pagination, err := client.Products.List(&gosalla.ListOptions{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}
	
	fmt.Printf("\nFound %d products (page %d of %d):\n\n", 
		len(products), pagination.CurrentPage, pagination.LastPage)
	
	for i, product := range products {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, product.Name, product.ID)
		fmt.Printf("   Price: %.2f, SKU: %s, Status: %s\n", 
			product.Price, product.SKU, product.Status)
		fmt.Println()
	}
	
	// Get a specific product by ID
	if len(products) > 0 {
		productID := products[0].ID
		fmt.Printf("\nFetching product ID %d...\n", productID)
		
		product, err := client.Products.Get(productID)
		if err != nil {
			log.Fatalf("Failed to get product: %v", err)
		}
		
		fmt.Printf("\nProduct Details:\n")
		fmt.Printf("Name: %s\n", product.Name)
		fmt.Printf("Description: %s\n", product.Description)
		fmt.Printf("Price: %.2f\n", product.Price)
		fmt.Printf("Quantity: %d\n", product.Quantity)
	}
	
	// Create a new product
	fmt.Println("\n\nCreating a new product...")
	newProduct := &gosalla.CreateProductRequest{
		Name:        "Test Product",
		Description: "This is a test product created via the Go SDK",
		Price:       99.99,
		Quantity:    100,
		SKU:         "TEST-SKU-001",
		Status:      "active",
	}
	
	created, err := client.Products.Create(newProduct)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	
	fmt.Printf("Successfully created product: %s (ID: %d)\n", created.Name, created.ID)
	
	// Update the product
	fmt.Println("\nUpdating the product...")
	updateReq := &gosalla.UpdateProductRequest{
		Price: 79.99,
	}
	
	updated, err := client.Products.Update(created.ID, updateReq)
	if err != nil {
		log.Fatalf("Failed to update product: %v", err)
	}
	
	fmt.Printf("Successfully updated product price to: %.2f\n", updated.Price)
	
	// Delete the product
	fmt.Println("\nDeleting the product...")
	if err := client.Products.Delete(created.ID); err != nil {
		log.Fatalf("Failed to delete product: %v", err)
	}
	
	fmt.Println("Successfully deleted the product")
}
