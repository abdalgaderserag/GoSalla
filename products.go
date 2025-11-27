package gosalla

import (
	"fmt"
	"time"
)

// ProductsService handles communication with the product-related endpoints
type ProductsService struct {
	client *Client
}

// Product represents a Salla product
type Product struct {
	ID              int                    `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	Price           float64                `json:"price"`
	SalePrice       float64                `json:"sale_price,omitempty"`
	SKU             string                 `json:"sku,omitempty"`
	Quantity        int                    `json:"quantity"`
	Status          string                 `json:"status"`
	Type            string                 `json:"type,omitempty"`
	Weight          float64                `json:"weight,omitempty"`
	CategoryID      int                    `json:"category_id,omitempty"`
	BrandID         int                    `json:"brand_id,omitempty"`
	Images          []ProductImage         `json:"images,omitempty"`
	Options         []ProductOption        `json:"options,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at,omitempty"`
	UpdatedAt       time.Time              `json:"updated_at,omitempty"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Alt      string `json:"alt,omitempty"`
	Position int    `json:"position,omitempty"`
}

// ProductOption represents a product option (e.g., size, color)
type ProductOption struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Values   []string `json:"values,omitempty"`
	Required bool     `json:"required,omitempty"`
}

// ProductsListResponse represents the response from listing products
type ProductsListResponse struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Data       []Product   `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// ProductResponse represents the response for a single product
type ProductResponse struct {
	Success bool    `json:"success"`
	Code    int     `json:"code"`
	Data    Product `json:"data"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price"`
	SalePrice   float64                `json:"sale_price,omitempty"`
	SKU         string                 `json:"sku,omitempty"`
	Quantity    int                    `json:"quantity"`
	Status      string                 `json:"status,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Weight      float64                `json:"weight,omitempty"`
	CategoryID  int                    `json:"category_id,omitempty"`
	BrandID     int                    `json:"brand_id,omitempty"`
	Images      []string               `json:"images,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price,omitempty"`
	SalePrice   float64                `json:"sale_price,omitempty"`
	SKU         string                 `json:"sku,omitempty"`
	Quantity    int                    `json:"quantity,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Weight      float64                `json:"weight,omitempty"`
	CategoryID  int                    `json:"category_id,omitempty"`
	BrandID     int                    `json:"brand_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// List retrieves all products with optional pagination
func (s *ProductsService) List(opts *ListOptions) ([]Product, *Pagination, error) {
	path := "/products"
	
	// Add query parameters
	if opts != nil {
		path += fmt.Sprintf("?page=%d&per_page=%d", opts.Page, opts.PerPage)
	}
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var resp ProductsListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, nil, err
	}
	
	return resp.Data, resp.Pagination, nil
}

// Get retrieves a product by ID
func (s *ProductsService) Get(id int) (*Product, error) {
	path := fmt.Sprintf("/products/%d", id)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp ProductResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// GetBySKU retrieves a product by SKU
func (s *ProductsService) GetBySKU(sku string) (*Product, error) {
	path := fmt.Sprintf("/products/sku/%s", sku)
	
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var resp ProductResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Create creates a new product
func (s *ProductsService) Create(product *CreateProductRequest) (*Product, error) {
	path := "/products"
	
	req, err := s.client.newRequest("POST", path, product)
	if err != nil {
		return nil, err
	}
	
	var resp ProductResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Update updates an existing product
func (s *ProductsService) Update(id int, product *UpdateProductRequest) (*Product, error) {
	path := fmt.Sprintf("/products/%d", id)
	
	req, err := s.client.newRequest("PUT", path, product)
	if err != nil {
		return nil, err
	}
	
	var resp ProductResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	
	return &resp.Data, nil
}

// Delete deletes a product
func (s *ProductsService) Delete(id int) error {
	path := fmt.Sprintf("/products/%d", id)
	
	req, err := s.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	
	return s.client.do(req, nil)
}

// ChangeStatus changes the status of a product
func (s *ProductsService) ChangeStatus(id int, status string) error {
	path := fmt.Sprintf("/products/%d/status", id)
	
	body := map[string]string{"status": status}
	req, err := s.client.newRequest("POST", path, body)
	if err != nil {
		return err
	}
	
	return s.client.do(req, nil)
}
